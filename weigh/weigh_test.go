package weigh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNeatSize(t *testing.T) {
	cases := []struct {
		expected string
		size     int64
	}{
		{"123 bytes", 123},
		{"12.06 KiB", 12345},
		{"11.77 MiB", 12345678},
		{"11.50 GiB", 12345678901},
		{"11.23 TiB", 12345678012345},
		{"10.97 PiB", 12345678012345678},
	}

	for _, tc := range cases {
		require.Equal(t, tc.expected, neatSize(tc.size))
	}
}

func TestDirBytes(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), 1024)
	writeFile(t, filepath.Join(dir, "b.txt"), 2048)

	sub := filepath.Join(dir, "sub")
	require.NoError(t, os.Mkdir(sub, 0o755))
	writeFile(t, filepath.Join(sub, "c.txt"), 512)

	got := dirBytes(dir)
	require.Equal(t, int64(1024+2048+512), got)
}

func TestDirBytes_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	require.Equal(t, int64(0), dirBytes(dir))
}

func TestDirBytes_MissingDir(t *testing.T) {
	// Should return 0 rather than panic.
	got := dirBytes(filepath.Join(t.TempDir(), "nonexistent"))
	require.Equal(t, int64(0), got)
}

func TestTopDir(t *testing.T) {
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "file.txt"), 100)

	sub := filepath.Join(dir, "subdir")
	require.NoError(t, os.Mkdir(sub, 0o755))
	writeFile(t, filepath.Join(sub, "inner.txt"), 200)

	results := topDir(dir)
	require.Len(t, results, 2)

	byName := make(map[string]SummaryData)
	for _, r := range results {
		byName[r.Name] = r
	}

	file := byName[filepath.Join(dir, "file.txt")]
	require.Equal(t, int64(100), file.Bytes)
	require.False(t, file.IsDir)

	sd := byName[sub]
	require.Equal(t, int64(200), sd.Bytes)
	require.True(t, sd.IsDir)
}

func TestTopDir_EmptyDir(t *testing.T) {
	results := topDir(t.TempDir())
	require.Empty(t, results)
}

func TestSummarize_DefaultsToCurrentDir(t *testing.T) {
	w := Weigh{}
	w.Summarize()
	// After Summarize with no paths, Paths must be set and Summaries populated
	// (the current directory always has at least the package source files).
	require.NotEmpty(t, w.Paths)
}

func TestSummarize_SingleTarget(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "x.txt"), 500)
	sub := filepath.Join(dir, "d")
	require.NoError(t, os.Mkdir(sub, 0o755))
	writeFile(t, filepath.Join(sub, "y.txt"), 300)

	w := Weigh{Paths: []string{dir}}
	w.Summarize()

	require.NotEmpty(t, w.Summaries)

	var total int64
	for _, s := range w.Summaries {
		total += s.Bytes
	}
	require.Equal(t, int64(500+300), total)
}

// writeFile creates a file of exactly size bytes filled with zeros.
func writeFile(t *testing.T, path string, size int) {
	t.Helper()
	f, err := os.Create(path)
	require.NoError(t, err)
	defer f.Close()
	require.NoError(t, f.Truncate(int64(size)))
}

func BenchmarkWeigh(b *testing.B) {
	for b.Loop() {
		w := Weigh{Paths: []string{"/var"}}
		w.Summarize()
	}
}
