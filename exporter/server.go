package exporter

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/xaque208/weigh/weigh"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	weighDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "weigh_duration_seconds",
			Help: "Duration of collections by the weigh",
		},
		[]string{"path"},
	)
	weighTargetGuage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "weigh",
			Name:      "target",
			Help:      "Weigh target",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(
		weighDuration,
		weighTargetGuage,
	)
}

func handler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	if target == "" {
		http.Error(w, "'target' parameter must be specified", 400)
		return
	}

	start := time.Now()
	registry := prometheus.NewRegistry()

	x := weigh.Weigh{Paths: []string{target}}
	x.Summarize()
	for _, s := range x.Summaries {
		weighTargetGuage.With(prometheus.Labels{"path": s.Name}).Set(float64(s.Bytes))
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
	duration := time.Since(start).Seconds()
	weighDuration.WithLabelValues("weigh").Observe(duration)
}

func StartMetricsServer(bindAddr string) {
	d := http.NewServeMux()

	d.Handle("/metrics", promhttp.Handler())
	d.HandleFunc("/weigh", handler)

	err := http.ListenAndServe(bindAddr, d)
	if err != nil {
		log.Fatal("Failed to start metrics server, error is:", err)
	}
}
