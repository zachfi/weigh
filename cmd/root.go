package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xaque208/weigh/exporter"
	"github.com/xaque208/weigh/weigh"
)

var rootCmd = &cobra.Command{
	Use:   "weigh",
	Short: "A thing to count bytes on files ",
	Long:  "",
	Args:  cobra.MinimumNArgs(0),
	Run:   run,
}

var (
	verbose       bool
	listenAddress string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Increase verbosity")
	rootCmd.PersistentFlags().StringVarP(&listenAddress, "listen", "L", "", "The listen address (default is not to listen (:9100)")
}

func run(cmd *cobra.Command, args []string) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if listenAddress != "" {
		exporter.StartMetricsServer(listenAddress)
	} else {
		w := weigh.Weigh{Paths: args}
		w.Summarize()
		w.Report()
	}
}
