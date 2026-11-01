package nvelope

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/muir/nject/v2"

	"github.com/pkg/errors"
)

// LogFlusher is used to check if a logger implements
// Flush().  This is useful as part of a panic handler.
type LogFlusher interface {
	Flush()
}

type panicError struct {
	msg   string
	r     interface{}
	stack string
}

func (err panicError) Error() string {
	return "panic: " + err.msg
}

// SetErrorOnPanic should be called as a defer.  It
// sets an error value if there is a panic.
//
// A panic of http.ErrAbortHandler is raised again.  net/http, including
// httputil.ReverseProxy when a streaming copy fails, panics with that
// sentinel so the server can abort the response without logging a stack
// trace.  The comparison is by identity, matching net/http.  A wrapped
// error is recorded like any other panic.
func SetErrorOnPanic(ep *error, log BasicLogger) {
	r := recover()
	if r == nil {
		return
	}
	if r == http.ErrAbortHandler {
		panic(r)
	}
	pe := panicError{
		msg:   fmt.Sprint(r),
		r:     r,
		stack: string(debug.Stack()),
	}
	*ep = errors.WithStack(pe)
	log.Error("panic!", map[string]interface{}{
		"msg":   pe.msg,
		"stack": pe.stack,
	})
	if flusher, ok := log.(LogFlusher); ok {
		flusher.Flush()
	}
}

// CatchPanic is a wrapper that catches downstream panics and returns
// an error a downsteam provider panic's.  http.ErrAbortHandler is
// raised again; see SetErrorOnPanic.
var CatchPanic = nject.Provide("catch-panic", catchPanicInjector)

func catchPanicInjector(inner func() error, log BasicLogger) (err error) {
	defer SetErrorOnPanic(&err, log)
	err = inner()
	return
}

// RecoverInterface returns the interface{} that recover()
// originally provided.  Or it returns nil if the
// error isn't a from a panic recovery.  This works only
// in conjunction with SetErrorOnPanic() and CatchPanic.
func RecoverInterface(err error) interface{} {
	if pe, ok := isPanicError(err); ok {
		return pe.r
	}
	return nil
}

// RecoverStack returns the stack from when recover()
// originally caught the panic.  Or it returns "" if the
// error isn't a from a panic recovery.  This works only
// in conjunction with SetErrorOnPanic() and CatchPanic.
func RecoverStack(err error) string {
	if pe, ok := isPanicError(err); ok {
		return pe.stack
	}
	return ""
}

func isPanicError(err error) (panicError, bool) {
	var pe panicError
	if errors.As(err, &pe) {
		return pe, true
	}
	return panicError{}, false
}
