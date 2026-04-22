package exporter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

func newTestExporter(t *testing.T, targets []string) *Exporter {
	t.Helper()
	e, err := New(Config{Targets: targets}, log.NewNopLogger(), nil)
	require.NoError(t, err)
	return e
}

func collectMetrics(e *Exporter) []*dto.MetricFamily {
	reg := prometheus.NewRegistry()
	reg.MustRegister(e)
	mfs, err := reg.Gather()
	if err != nil {
		// Partial gather is still useful in tests.
		_ = err
	}
	return mfs
}

func TestCollect_ValidTarget(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), 1024)
	sub := filepath.Join(dir, "sub")
	require.NoError(t, os.Mkdir(sub, 0o755))
	writeFile(t, filepath.Join(sub, "b.txt"), 512)

	e := newTestExporter(t, []string{dir})
	mfs := collectMetrics(e)

	byName := metricFamiliesByName(mfs)

	require.Contains(t, byName, "weigh_duration_seconds")
	require.Contains(t, byName, "weigh_target_bytes")
	require.Contains(t, byName, "weigh_collect_errors_total")

	// One duration sample per target.
	require.Len(t, byName["weigh_duration_seconds"].GetMetric(), 1)
	require.Equal(t, dir, byName["weigh_duration_seconds"].GetMetric()[0].GetLabel()[0].GetValue())

	// Error counter should be zero.
	errVal := byName["weigh_collect_errors_total"].GetMetric()[0].GetGauge().GetValue()
	require.Equal(t, float64(0), errVal)

	// Bytes for the subdirectory should equal 512.
	var subdirBytes float64
	for _, m := range byName["weigh_target_bytes"].GetMetric() {
		for _, lp := range m.GetLabel() {
			if lp.GetValue() == sub {
				subdirBytes = m.GetGauge().GetValue()
			}
		}
	}
	require.Equal(t, float64(512), subdirBytes)
}

func TestCollect_NoTargets(t *testing.T) {
	e := newTestExporter(t, nil)
	mfs := collectMetrics(e)

	byName := metricFamiliesByName(mfs)
	// Only the error counter should be present with value 0.
	require.Contains(t, byName, "weigh_collect_errors_total")
	errVal := byName["weigh_collect_errors_total"].GetMetric()[0].GetGauge().GetValue()
	require.Equal(t, float64(0), errVal)
}

func TestCollect_PathCleaning(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "f.txt"), 256)

	// Pass a path with trailing slash and redundant components; label should be the cleaned absolute path.
	dirty := dir + "/sub/../"
	e := newTestExporter(t, []string{dirty})
	mfs := collectMetrics(e)

	byName := metricFamiliesByName(mfs)
	require.Contains(t, byName, "weigh_duration_seconds")

	label := byName["weigh_duration_seconds"].GetMetric()[0].GetLabel()[0].GetValue()
	require.Equal(t, dir, label, "duration label should be cleaned absolute path")
}

func metricFamiliesByName(mfs []*dto.MetricFamily) map[string]*dto.MetricFamily {
	m := make(map[string]*dto.MetricFamily, len(mfs))
	for _, mf := range mfs {
		m[mf.GetName()] = mf
	}
	return m
}

func writeFile(t *testing.T, path string, size int) {
	t.Helper()
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, f.Truncate(int64(size)))
}
