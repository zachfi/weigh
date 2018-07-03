package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

import log "github.com/sirupsen/logrus"

const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
	PETABYTE = 1024 * TERABYTE
)

type summaryData struct {
	Name  string
	Bytes int64
}

type summariesData []summaryData

func (slice summariesData) Len() int {
	return len(slice)
}

func (slice summariesData) Less(i, j int) bool {
	return slice[i].Bytes < slice[j].Bytes
}

func (slice summariesData) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

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
	log.Infof("Entering directory %s", directory)
	var dirSize int64 = 0

	countDir := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dirSize += info.Size()
		}

		if err != nil {
			fmt.Print(err)
		}

		return nil
	}

	filepath.Walk(directory, countDir)

	return dirSize
}

func report(summaries summariesData) {
	var total int64 = 0

	sort.Sort(summaries)

	for _, item := range summaries {
		if item.Bytes == 0 {
			continue
		}

		total += item.Bytes

		fi, err := os.Stat(item.Name)

		switch {
		case err != nil:
			log.Error(err)
		case fi.IsDir():
			fmt.Printf("%15s    %s/\n", neatSize(item.Bytes), item.Name)
		default:
			fmt.Printf("%15s    %s\n", neatSize(item.Bytes), item.Name)
		}

	}

	fmt.Printf("%16s %s\n", "---", "---")
	fmt.Printf("%15s  %s\n", neatSize(total), ":total size")
	fmt.Printf("%16s %s\n", "---", "---")

}

func topDir(directory string) summariesData {
	summary := summariesData{}

	files, _ := ioutil.ReadDir(directory)

	for _, f := range files {
		fullpath := filepath.Join(directory, f.Name())

		if f.IsDir() {
			summary = append(summary, summaryData{Name: fullpath, Bytes: dirBytes(fullpath)})
		} else {
			summary = append(summary, summaryData{Name: fullpath, Bytes: f.Size()})
		}
	}

	return summary
}

func main() {

	var verbose bool
	flag.BoolVar(&verbose, "v", false, "Increase verbosity")

	flag.Parse()

	if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	directories := flag.Args()

	if len(directories) == 0 {
		directories = append(directories, "./")
	}

	// summaries := make([]summaryData, 0)
	summaries := summariesData{}

	for _, d := range directories {
		for _, sum := range topDir(d) {
			summaries = append(summaries, sum)
		}
	}

	report(summaries)
}
