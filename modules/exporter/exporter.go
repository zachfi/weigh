package exporter

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"github.com/zachfi/weigh/weigh"

	"github.com/grafana/dskit/services"
)

var (
	weighDurationDesc = prometheus.NewDesc(
		"weigh_duration_seconds",
		"Time taken to weigh the targets",
		[]string{"path"},
		nil,
	)
	weighTargetDesc = prometheus.NewDesc(
		"weigh_target_bytes",
		"Weigh target bytes",
		[]string{"path"},
		nil,
	)
)

type Exporter struct {
	services.Service
	cfg *Config

	logger log.Logger
}

func New(cfg Config, logger log.Logger, conn *grpc.ClientConn) (*Exporter, error) {
	logger = log.With(logger, "module", "exporter")

	e := &Exporter{
		cfg:    &cfg,
		logger: logger,
	}

	e.Service = services.NewBasicService(e.starting, e.running, e.stopping)
	return e, nil
}

func (e *Exporter) starting(ctx context.Context) error {
	return nil
}

func (e *Exporter) running(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func (e *Exporter) stopping(_ error) error {
	return nil
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- weighDurationDesc
	ch <- weighTargetDesc
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, target := range e.cfg.Targets {
		x := weigh.Weigh{Paths: []string{target}}

		t := time.Now()
		x.Summarize()
		ch <- prometheus.MustNewConstMetric(weighDurationDesc, prometheus.GaugeValue, float64(time.Since(t).Seconds()), target)

		for _, s := range x.Summaries {
			ch <- prometheus.MustNewConstMetric(weighTargetDesc, prometheus.GaugeValue, float64(s.Bytes), s.Name)
		}
	}
}
