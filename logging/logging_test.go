package logging

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"testing"
)

func TestNewModuleLog(t *testing.T) {
	log := NewModuleLog("test")

	assert.IsType(t, &zap.SugaredLogger{}, log)
}
