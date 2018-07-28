package test_helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestWithCaptureLogsJSON(t *testing.T) {
	assert.IsType(t, &CaptureLogsJSON{}, WithCaptureLogsJSON())
}

func TestCaptureLogsJSON_Before(t *testing.T) {
	clj := WithCaptureLogsJSON()

	ctx := &Context{}

	clj.Before(ctx)

	zap.S().Debug("this is a test")

	logs := ctx.Log.String()

	t.Log(logs)

	assert.Contains(t, logs, "this is a test")
}

func TestCaptureLogsJSON_After(t *testing.T) {
	clj := WithCaptureLogsJSON()

	ctx := &Context{}

	clj.After(ctx)

	zap.S().Debug("this is a test")

	logs := SuiteLog.String()

	t.Log(logs)

	assert.Contains(t, logs, "this is a test")
}
