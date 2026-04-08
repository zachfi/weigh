package weigh

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"

	log "github.com/sirupsen/logrus"
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
		log.Debugf("Adding default path")
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
	if bytes >= PETABYTE {
		return fmt.Sprintf("%.2f PiB", float64(bytes)/float64(PETABYTE))
	}

	if bytes >= TERABYTE {
		return fmt.Sprintf("%.2f TiB", float64(bytes)/float64(TERABYTE))
	}

	if bytes >= GIGABYTE {
		return fmt.Sprintf("%.2f GiB", float64(bytes)/float64(GIGABYTE))
	}

	if bytes >= MEGABYTE {
		return fmt.Sprintf("%.2f MiB", float64(bytes)/float64(MEGABYTE))
	}

	if bytes >= KILOBYTE {
		return fmt.Sprintf("%.2f KiB", float64(bytes)/float64(KILOBYTE))
	}

	return fmt.Sprintf("%d bytes", int64(bytes))
}

func dirBytes(directory string) int64 {
	log.Debugf("Entering directory %s", directory)
	var dirSize int64 = 0

	countDir := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error(err)
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				log.Error(err)
				return nil
			}
			dirSize += info.Size()
		}
		return nil
	}

	if err := filepath.WalkDir(directory, countDir); err != nil {
		log.Error(err)
	}

	return dirSize
}

func topDir(directory string) SummariesData {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Error(err)
		return SummariesData{}
	}

	results := make(SummariesData, len(files))
	var wg sync.WaitGroup

	for i, f := range files {
		fullpath := filepath.Join(directory, f.Name())

		if f.IsDir() {
			wg.Add(1)
			go func(idx int, path string) {
				defer wg.Done()
				results[idx] = SummaryData{Name: path, Bytes: dirBytes(path), IsDir: true}
			}(i, fullpath)
		} else {
			info, err := f.Info()
			if err != nil {
				log.Error(err)
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
