package logging

import (
	"log"
	"strings"

	"github.com/bgetsug/go-toolbox/config"
	"go.uber.org/zap"
)

func Init(env config.Environment) {
	logger, err := zap.NewProduction()

	if err != nil {
		log.Fatal(err)
	}

	if env == config.LOCAL || env == config.TESTING {
		logger, err = zap.NewDevelopment()

		if err != nil {
			log.Fatal(err)
		}
	}

	zap.ReplaceGlobals(logger)
}

// NewModuleLog return a zap SugaredLogger configured to log the
// specified name as the "module" field in the logging context
func NewModuleLog(name ...string) *zap.SugaredLogger {
	return zap.S().Named(strings.Join(name, "."))
}
