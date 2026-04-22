package weigh

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
)

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
	PETABYTE = 1024 * TERABYTE
)

type Weigh struct {
	Paths     []string
	Summaries SummariesData
}

func (w *Weigh) Summarize() {
	if len(w.Paths) == 0 {
		w.Paths = []string{"./"}
	}

	summaries := make(SummariesData, 0, len(w.Paths))

	for _, d := range w.Paths {
		summaries = append(summaries, topDir(d)...)
	}

	w.Summaries = summaries
}

func (w *Weigh) Report() {
	summaries := w.Summaries
	var total int64 = 0

	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Bytes > summaries[j].Bytes
	})

	sep := fmt.Sprintf("%15s    %s", "---", "---")
	for _, item := range summaries {
		if item.Bytes == 0 {
			continue
		}

		total += item.Bytes

		if item.IsDir {
			fmt.Printf("%15s    %s/\n", neatSize(item.Bytes), item.Name)
		} else {
			fmt.Printf("%15s    %s\n", neatSize(item.Bytes), item.Name)
		}
	}

	fmt.Println(sep)
	fmt.Printf("%15s    %s\n", neatSize(total), "total")
	fmt.Println(sep)
}

type SummaryData struct {
	Name  string
	Bytes int64
	IsDir bool
}

type SummariesData []SummaryData

func neatSize(bytes int64) string {
	switch {
	case bytes >= PETABYTE:
		return fmt.Sprintf("%.2f PiB", float64(bytes)/float64(PETABYTE))
	case bytes >= TERABYTE:
		return fmt.Sprintf("%.2f TiB", float64(bytes)/float64(TERABYTE))
	case bytes >= GIGABYTE:
		return fmt.Sprintf("%.2f GiB", float64(bytes)/float64(GIGABYTE))
	case bytes >= MEGABYTE:
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(MEGABYTE))
	case bytes >= KILOBYTE:
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(KILOBYTE))
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}

func dirBytes(directory string) int64 {
	var dirSize int64

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			slog.Debug("walk error", "path", path, "err", err)
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				slog.Debug("stat error", "path", path, "err", err)
				return nil
			}
			dirSize += info.Size()
		}
		return nil
	})
	if err != nil {
		slog.Error("walkdir failed", "dir", directory, "err", err)
	}

	return dirSize
}

// workers bounds the concurrent goroutines spawned by topDir.
var workers = runtime.NumCPU() * 2

func topDir(directory string) SummariesData {
	files, err := os.ReadDir(directory)
	if err != nil {
		slog.Error("readdir failed", "dir", directory, "err", err)
		return SummariesData{}
	}

	results := make(SummariesData, len(files))
	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	for i, f := range files {
		fullpath := filepath.Join(directory, f.Name())

		if f.IsDir() {
			wg.Add(1)
			sem <- struct{}{}
			go func(idx int, path string) {
				defer func() {
					<-sem
					wg.Done()
				}()
				results[idx] = SummaryData{Name: path, Bytes: dirBytes(path), IsDir: true}
			}(i, fullpath)
		} else {
			info, err := f.Info()
			if err != nil {
				slog.Debug("stat error", "path", fullpath, "err", err)
				continue
			}
			results[i] = SummaryData{Name: fullpath, Bytes: info.Size()}
		}
	}

	wg.Wait()

	summary := make(SummariesData, 0, len(results))
	for _, r := range results {
		if r.Name != "" {
			summary = append(summary, r)
		}
	}

	return summary
}
