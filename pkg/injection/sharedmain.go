package injection

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"knative.dev/pkg/logging"

	pkghandlers "github.com/itsmurugappan/http-handlers/pkg/handlers/favicon"
)

const (
	log_config = `{
					        "level": "info",
					        "development": false,
					        "outputPaths": ["stdout"],
					        "errorOutputPaths": ["stderr"],
					        "encoding": "json",
					        "encoderConfig": {
					          "timeKey": "ts",
					          "levelKey": "level",
					          "nameKey": "logger",
					          "callerKey": "caller",
					          "messageKey": "msg",
					          "stacktraceKey": "stacktrace",
					          "lineEnding": "",
					          "levelEncoder": "",
					          "timeEncoder": "iso8601",
					          "durationEncoder": "",
					          "callerEncoder": ""
					        }
      					}`
)

// SharedMain registers http handlers and creates other boiler plate
// for a main method to use
func SharedMain(ctx context.Context, port string, handlers ...func(context.Context)) {
	logger, _ := logging.NewLogger(log_config, "")
	ctx = logging.WithLogger(ctx, logger)

	//register handlers
	for _, h := range handlers {
		h(ctx)
	}
	//add generic handlers
	pkghandlers.HandleFavicon()

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}
