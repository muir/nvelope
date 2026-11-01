package nvelope_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muir/nject/v2"
	"github.com/muir/nvelope"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordingLogger struct {
	msgs []string
}

func (l *recordingLogger) Debug(msg string, fields ...map[string]interface{}) {}
func (l *recordingLogger) Warn(msg string, fields ...map[string]interface{})  {}
func (l *recordingLogger) Error(msg string, fields ...map[string]interface{}) {
	l.msgs = append(l.msgs, msg)
}

func TestSetErrorOnPanicRepanicsErrAbortHandler(t *testing.T) {
	err := assert.AnError
	log := &recordingLogger{}
	func() {
		defer func() {
			r := recover()
			assert.Equal(t, http.ErrAbortHandler, r)
		}()
		defer nvelope.SetErrorOnPanic(&err, log)
		panic(http.ErrAbortHandler)
	}()
	assert.Equal(t, assert.AnError, err)
	assert.Empty(t, log.msgs)
}

func TestSetErrorOnPanicRecoversOtherPanics(t *testing.T) {
	var err error
	log := &recordingLogger{}
	func() {
		defer nvelope.SetErrorOnPanic(&err, log)
		panic("boom")
	}()
	require.Error(t, err)
	assert.Equal(t, "panic: boom", err.Error())
	assert.Equal(t, "boom", nvelope.RecoverInterface(err))
	assert.NotEmpty(t, nvelope.RecoverStack(err))
	assert.Equal(t, []string{"panic!"}, log.msgs)

	// Same text as the sentinel, but not the sentinel itself.  net/http
	// only suppresses the exact value, so this stays a normal panic.
	err = nil
	log.msgs = nil
	lookalike := http.ErrAbortHandler.Error()
	func() {
		defer nvelope.SetErrorOnPanic(&err, log)
		panic(lookalike)
	}()
	require.EqualError(t, err, "panic: net/http: abort Handler")
	assert.Equal(t, []string{"panic!"}, log.msgs)
}

func TestCatchPanicLetsAbortHandlerPropagate(t *testing.T) {
	var invoke func() error
	err := nject.Sequence("abort",
		nvelope.NoLogger,
		nvelope.CatchPanic,
		func() error {
			panic(http.ErrAbortHandler)
		},
	).Bind(&invoke, nil)
	require.NoError(t, err)

	defer func() {
		assert.Equal(t, http.ErrAbortHandler, recover())
	}()
	got := invoke()
	t.Fatalf("invoke returned %v", got)
}

func TestCatchPanicConvertsOtherPanics(t *testing.T) {
	var invoke func() error
	err := nject.Sequence("panic",
		nvelope.NoLogger,
		nvelope.CatchPanic,
		func() error {
			panic("boom")
		},
	).Bind(&invoke, nil)
	require.NoError(t, err)

	got := invoke()
	require.EqualError(t, got, "panic: boom")
	assert.Equal(t, "boom", nvelope.RecoverInterface(got))
}

// abortEndpoint is the handler under test.  flushed copies a partial body
// into net/http and flushes that buffer onto the wire before the sentinel
// panic, which is the shape of a streaming handler that then aborts.
func abortEndpoint(flushed bool) func(*nvelope.DeferredWriter) (nvelope.Response, error) {
	return func(w *nvelope.DeferredWriter) (nvelope.Response, error) {
		if flushed {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("partial"))
			if err := w.Flush(); err != nil {
				return nil, err
			}
			// DeferredWriter.Flush only copies into net/http's buffer.
			// Flush the underlying writer so the partial body is on the wire.
			w.UnderlyingWriter().(http.Flusher).Flush()
		}
		panic(http.ErrAbortHandler)
	}
}

func TestErrAbortHandlerAbortsWithoutServerLog(t *testing.T) {
	for _, flushed := range []bool{false, true} {
		t.Run(map[bool]string{false: "before write", true: "after flush"}[flushed], func(t *testing.T) {
			var appLog bytes.Buffer
			var handler http.HandlerFunc
			err := nject.Sequence("abort-http",
				nvelope.LoggerFromStd(log.New(&appLog, "", 0)),
				nvelope.InjectWriter,
				nvelope.EncodeJSON,
				nvelope.CatchPanic,
				abortEndpoint(flushed),
			).Bind(&handler, nil)
			require.NoError(t, err)

			var serverLog bytes.Buffer
			srv := httptest.NewUnstartedServer(handler)
			srv.Config.ErrorLog = log.New(&serverLog, "", 0)
			srv.Start()
			t.Cleanup(srv.Close)

			res, err := srv.Client().Get(srv.URL)
			if !flushed {
				// Nothing was written, so the abort closes the connection
				// before a response exists.
				require.Error(t, err)
				if res != nil {
					_ = res.Body.Close()
				}
			} else {
				// Headers and a partial body reached the client. Reading
				// stops with unexpected EOF because the server closed the
				// connection instead of finishing the response.
				require.NoError(t, err)
				defer func() { _ = res.Body.Close() }()
				body, readErr := io.ReadAll(res.Body)
				assert.Equal(t, http.StatusOK, res.StatusCode)
				assert.Equal(t, "partial", string(body))
				assert.ErrorIs(t, readErr, io.ErrUnexpectedEOF)
			}

			assert.NotContains(t, serverLog.String(), "panic serving")
			assert.NotContains(t, appLog.String(), "panic!")
		})
	}
}
