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
func SharedMain(ctx context.Context, port, component string, handlerInits ...func(context.Context, configmap.Watcher) Impl) {
	logger, atomicLevel := sharedmain.SetupLoggerOrDie(ctx, component)
	defer flush(logger)
	ctx = logging.WithLogger(ctx, logger)
	ctx, _ = injection.Default.SetupInformers(ctx, sharedmain.ParseAndGetConfigOrDie())

	//start config map watchers
	cmw := informer.NewInformedWatcher(kubeclient.Get(ctx), system.Namespace())
	sharedmain.WatchLoggingConfigOrDie(ctx, cmw, logger, atomicLevel, component)

	//instantiate handlers and construct handler impl list
	impls := make([]Impl, len(handlerInits))
	for i, h := range handlerInits {
		impls[i] = h(ctx, cmw)
	}
	//add generic handlers
	pkghandlers.HandleFavicon()

	logger.Info("Starting configuration manager...")
	if err := cmw.Start(ctx.Done()); err != nil {
		logger.Fatalw("Failed to start configuration manager", err)
	}

	//start handler impls
	for _, i := range impls {
		i.Start(ctx)
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}

func flush(logger *zap.SugaredLogger) {
	logger.Sync()
}

// Impl is handler interfce
type Impl interface {
	Start(context.Context)
}
