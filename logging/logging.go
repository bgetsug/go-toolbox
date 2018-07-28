package logging

import (
	"log"
	"strings"

	"github.com/bgetsug/go-toolbox/config"
	"go.uber.org/zap"
)

var (
	// Logger is a global instance of a zap Logger
	Logger *zap.Logger

	// Log is a global instance of a zap SugaredLogger
	Log *zap.SugaredLogger
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

	Logger = logger
	Log = Logger.Sugar()
	zap.ReplaceGlobals(Logger)
}

// NewModuleLog return a zap SugaredLogger configured to log the
// specified name as the "module" field in the logging context
func NewModuleLog(name ...string) *zap.SugaredLogger {
	return Log.Named(strings.Join(name, "."))
}
