package test_helpers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/gavv/httpexpect.v1"
)

var (
	SuiteLog bytes.Buffer
)

// Decorator represents the behavioral contract of a test decorator, providing Before and After "hook" functionality.
type Decorator interface {
	Before(ctx *Context)
	After(ctx *Context)
}

// WithDecorators should be used to wrap a test function (f) with the desired decorators. Decorator Before and After
// methods will execute sequentially, before and after the test function (f), respectively.
func WithDecorators(t *testing.T, decorators []Decorator, f func(ctx *Context)) {
	ctx := &Context{
		Props: map[string]interface{}{
			"t": t,
		},
	}

	for _, dec := range decorators {
		dec.Before(ctx)
	}

	f(ctx) // test runs here

	for _, dec := range decorators {
		dec.After(ctx)
	}
}

type CaptureLogsJSON struct{}

func WithCaptureLogsJSON() *CaptureLogsJSON {
	return &CaptureLogsJSON{}
}

func (d *CaptureLogsJSON) Before(ctx *Context) {
	d.newZap(&ctx.Log)
}

func (d *CaptureLogsJSON) After(ctx *Context) {
	d.newZap(&SuiteLog)
}

func (d *CaptureLogsJSON) newZap(w io.Writer) {
	logger := zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
			zapcore.AddSync(w),
			zap.DebugLevel,
		),
		zap.Development(),
	)

	zap.ReplaceGlobals(logger)
}

// HTTPExpect ...
type HTTPExpect struct {
	*httpexpect.Expect
	handler http.Handler
	server  *httptest.Server
}

// WithHTTPExpect ...
func WithHTTPExpect(handler http.Handler) *HTTPExpect {
	return &HTTPExpect{handler: handler}
}

// Before ...
func (d *HTTPExpect) Before(ctx *Context) {
	t := ctx.MustGet("t").(*testing.T)

	d.server = httptest.NewServer(d.handler)

	d.Expect = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  d.server.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCurlPrinter(t),
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

// After ...
func (d *HTTPExpect) After(ctx *Context) {
	d.server.Close()
}
