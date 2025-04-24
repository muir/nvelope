package nvelope_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/muir/nvelope"

	"github.com/muir/nject/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type logWrap struct {
	t *testing.T
}

func (l logWrap) Print(v ...any) {
	l.t.Log(v...)
}

func logFromT(t *testing.T) func() nvelope.BasicLogger {
	return nvelope.LoggerFromStd(logWrap{t: t})
}

func doTest(t *testing.T, chain ...any) ([]byte, *http.Response) {
	var handler func(http.ResponseWriter, *http.Request)
	err := nject.Sequence("test",
		logFromT(t),
		nvelope.InjectWriter,
		nvelope.AutoFlushWriter,
		nvelope.EncodeJSON,
		nject.Sequence("chain", chain...),
	).Bind(&handler, nil)
	require.NoError(t, err, nject.DetailedError(err))

	ts := httptest.NewServer(http.HandlerFunc(handler))
	resp, err := http.Get(ts.URL + "/irrelevant")
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return body, resp
}

func TestJSONEncoderJSONString(t *testing.T) {
	body, resp := doTest(t,
		func() (nvelope.Response, error) {
			return "foo", nil
		})
	require.Equal(t, `"foo"`, string(body))
	require.Equal(t, 200, resp.StatusCode)
}

func TestJSONEncoderJSONStruct(t *testing.T) {
	body, resp := doTest(t,
		func() (nvelope.Response, error) {
			return struct{ Foo string }{Foo: "bar"}, nil
		})
	require.Equal(t, `{"Foo":"bar"}`, string(body))
	require.Equal(t, 200, resp.StatusCode)
}

func TestJSONEncoderWriterDone(t *testing.T) {
	body, resp := doTest(t,
		func(w *nvelope.DeferredWriter) (nvelope.Response, error) {
			http.Error(w, "never mind", http.StatusUnauthorized)
			assert.NoError(t, w.Flush())
			return "foo", nil
		})
	require.Equal(t, "never mind\n", string(body))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJSONEncoderWriterUsed(t *testing.T) {
	body, resp := doTest(t,
		func(w *nvelope.DeferredWriter) (nvelope.Response, error) {
			http.Error(w, "never mind", http.StatusUnauthorized)
			return "foo", nil
		})
	require.Equal(t, "never mind\n", string(body))
	require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestJSONEncoderError(t *testing.T) {
	body, resp := doTest(t,
		func() (nvelope.Response, error) {
			return "", fmt.Errorf("howdy")
		})
	require.Equal(t, `howdy`, string(body))
	require.Equal(t, 500, resp.StatusCode)
}
