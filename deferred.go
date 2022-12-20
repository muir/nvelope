package nvelope

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// DeferredWriter that wraps an underlying http.ResponseWriter.
// DeferredWriter buffers writes and headers.  The buffer can be
// reset.  When it's time to actually write, use Flush().
type DeferredWriter struct {
	base        http.ResponseWriter
	passthrough bool
	header      http.Header
	buffer      []byte
	status      int
	resetHeader http.Header
	flushed     bool
}

// NewDeferredWriter returns a DeferredWriter based on a
// base ResponseWriter.  It re-injects the base writer
// so that in effect, there is only one writer present.
func NewDeferredWriter(w http.ResponseWriter) (*DeferredWriter, http.ResponseWriter) {
	dw := &DeferredWriter{
		base:        w,
		header:      w.Header().Clone(),
		resetHeader: w.Header().Clone(),
		buffer:      make([]byte, 0, 4*1024),
	}
	return dw, dw
}

// Header is the same as http.ResponseWriter.Header
func (w *DeferredWriter) Header() http.Header {
	if w.passthrough {
		return w.base.Header()
	}
	return w.header
}

// Write is the same as http.ResponseWriter.Write
// except that the action is delayed until Flush() is called.
func (w *DeferredWriter) Write(b []byte) (int, error) {
	if w.passthrough {
		return w.base.Write(b)
	}
	w.buffer = append(w.buffer, b...)
	return len(b), nil
}

// WriteHeader is the same as http.ResponseWriter.WriteHeader
// except that the action is delayed until Flush() is called.
func (w *DeferredWriter) WriteHeader(statusCode int) {
	if w.passthrough {
		w.base.WriteHeader(statusCode)
	} else {
		w.status = statusCode
	}
}

// Reset empties the DeferredWriter's buffers and resets its Header
// back to its original state.  Reset returns error if UnderlyingWriter()
// or Flush() have been called.
func (w *DeferredWriter) Reset() error {
	if w.passthrough {
		return errors.New("Attempt to reset a DeferredWriter after it is in passthrough mode")
	}
	w.buffer = nil
	w.status = 0
	w.header = w.resetHeader.Clone()
	return nil
}

// PreserveHeader saves the current Header so that a Reset will revert
// back to the header just saved.
func (w *DeferredWriter) PreserveHeader() {
	w.resetHeader = w.header.Clone()
}

// UnderlyingWriter returns the underlying writer.  Any header
// modifications made with the DeferredWriter are copied to the
// base writer.  After a call to UnderlyingWriter, the DeferredWriter
// switches to passthrough mode: all future calls to Write(),
// Header(), etc are passed through to the http.ResponseWriter that
// was used to initialize the DeferredWrited.
func (w *DeferredWriter) UnderlyingWriter() http.ResponseWriter {
	if w.passthrough {
		return w.base
	}
	w.passthrough = true
	h := w.base.Header()
	for k := range h {
		if v, ok := w.header[k]; ok {
			h[k] = v
		} else {
			delete(h, k)
		}
	}
	for k, v := range w.header {
		if _, ok := h[k]; ok {
			continue
		}
		h[k] = v
	}
	return w.base
}

// Flush pushes the buffered write content through to the base writer.
// You can only flush once.  After a flush, all further calls are passed
// through to be base writer.  WriteHeader() will be called on the base
// writer even if there is no buffered data.
func (w *DeferredWriter) Flush() error {
	if w.passthrough {
		return errors.New("Attempt flush deferred writer that is not deferred")
	}
	w.flushed = true
	base := w.UnderlyingWriter()
	if w.status != 0 {
		base.WriteHeader(w.status)
	}
	for i := 0; i < len(w.buffer)-1; {
		amt, err := base.Write(w.buffer[i:])
		if err != nil {
			// Is this handling of short writes necessary?  Perhaps
			// so since a follow-up write will probably give a
			// more accurate error.
			if errors.Is(err, io.ErrShortWrite) {
				i += amt
				continue
			}
			return errors.Wrap(err, "flush buffered writer")
		}
		break
	}
	return nil
}

// FlushIfNotFlushed calls Flush if the DeferredWriter is not in
// passthrough mode.
func (w *DeferredWriter) FlushIfNotFlushed() error {
	if !w.passthrough {
		return w.Flush()
	}
	return nil
}

// Done returns true if the DeferredWriter is in passthrough mode.
func (w *DeferredWriter) Done() bool {
	return w.passthrough
}

// Body returns the internal buffer used by DeferredWriter.  Do not modify it.
// It also returns the status code (if set).
// If UnderlyingWriter() has been called, then Body() will return an error since
// the underlying buffer does not represent what has been written.
func (w *DeferredWriter) Body() ([]byte, int, error) {
	if w.passthrough && !w.flushed {
		return nil, 0, errors.New("unable to provide body because DeferredWriter is operating in passthrough mode")
	}
	return w.buffer, w.status, nil
}
