package app

import (
	"context"
	"fmt"

	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/modules"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/weaveworks/common/server"

	"github.com/zachfi/weigh/modules/exporter"
)

const (
	Server string = "server"

	Once     string = "once"
	Exporter string = "exporter"
)

func (a *App) setupModuleManager() error {
	mm := modules.NewManager(a.logger)
	mm.RegisterModule(Server, a.initServer, modules.UserInvisibleModule)
	mm.RegisterModule(Once, nil)
	mm.RegisterModule(Exporter, a.initExporter)

	deps := map[string][]string{
		// Server:       nil,
		Exporter: {Server},
		Once:     {},
	}

	for mod, targets := range deps {
		if err := mm.AddDependency(mod, targets...); err != nil {
			return err
		}
	}

	a.ModuleManager = mm

	return nil
}

func (a *App) initExporter() (services.Service, error) {
	e, err := exporter.New(a.cfg.Exporter, a.logger, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to init "+metricsNamespace)
	}

	prometheus.MustRegister(e)

	return e, nil
}

func (a *App) initServer() (services.Service, error) {
	a.cfg.Server.MetricsNamespace = metricsNamespace
	a.cfg.Server.ExcludeRequestInLog = true
	a.cfg.Server.RegisterInstrumentation = true

	server, err := server.New(a.cfg.Server)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create server")
	}

	servicesToWaitFor := func() []services.Service {
		svs := []services.Service(nil)
		for m, s := range a.serviceMap {
			// Server should not wait for itself.
			if m != Server {
				svs = append(svs, s)
			}
		}

		return svs
	}

	a.Server = server

	serverDone := make(chan error, 1)

	runFn := func(ctx context.Context) error {
		go func() {
			defer close(serverDone)
			serverDone <- server.Run()
		}()

		select {
		case <-ctx.Done():
			return nil
		case err := <-serverDone:
			if err != nil {
				return err
			}

			return fmt.Errorf("server stopped unexpectedly")
		}
	}

	stoppingFn := func(_ error) error {
		// wait until all modules are done, and then shutdown server.
		for _, s := range servicesToWaitFor() {
			_ = s.AwaitTerminated(context.Background())
		}

		// shutdown HTTP and gRPC servers (this also unblocks Run)
		server.Shutdown()

		// if not closed yet, wait until server stops.
		<-serverDone
		_ = level.Info(a.logger).Log("msg", "server stopped")
		return nil
	}

	return services.NewBasicService(nil, runFn, stoppingFn), nil
}
