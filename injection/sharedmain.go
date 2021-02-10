package injection

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/configmap/informer"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/system"

	pkghandlers "github.com/itsmurugappan/http-handlers/handlers/favicon"
)

// SharedMain registers http handlers and creates other boiler plate
// for a main method to use
func SharedMain(ctx context.Context, port, component string, newStore func(configmap.Logger) *configmap.UntypedStore, handlers ...func(context.Context)) {
	logger, atomicLevel := sharedmain.SetupLoggerOrDie(ctx, component)
	defer flush(logger)
	ctx = logging.WithLogger(ctx, logger)
	ctx, _ = injection.Default.SetupInformers(ctx, sharedmain.ParseAndGetConfigOrDie())

	//start config map watchers
	store := newStore(logger)
	cmw := informer.NewInformedWatcher(kubeclient.Get(ctx), system.Namespace())
	store.WatchConfigs(cmw)
	sharedmain.WatchLoggingConfigOrDie(ctx, cmw, logger, atomicLevel, component)
	logger.Info("Starting configuration manager...")
	if err := cmw.Start(ctx.Done()); err != nil {
		logger.Fatalw("Failed to start configuration manager", err)
	}

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

func flush(logger *zap.SugaredLogger) {
	logger.Sync()
}
