package exporter

import (
	"encoding/json"
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

type HealthStatus struct{}

func init() {
	prometheus.MustRegister(
		weighDuration,
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
	registry.MustRegister(weighTargetGuage)

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

func statusHandler(w http.ResponseWriter, r *http.Request) {

	health := &HealthStatus{}

	bytes, err := json.MarshalIndent(health, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	log.Error(err)
}

func StartMetricsServer(bindAddr string) {
	d := http.NewServeMux()

	d.Handle("/metrics", promhttp.Handler())
	d.HandleFunc("/weigh", handler)
	d.HandleFunc("/status/check", statusHandler)

	err := http.ListenAndServe(bindAddr, d)
	if err != nil {
		log.Fatal("Failed to start metrics server, error is:", err)
	}
}
