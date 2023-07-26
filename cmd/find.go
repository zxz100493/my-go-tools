package cmd

import (
	"my-go-tools/service"

	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find a resource",
	Long:  "Find a resource by name or ID. For example:",
	Run:   findFile,
}

var (
	Verbose     bool
	Target      string
	Destination string
)

func init() {
	rootCmd.AddCommand(findCmd)
	parseArgs()
}

func findFile(cmd *cobra.Command, args []string) {
	cfg := &service.Config{
		Visible:     Verbose,
		Target:      Target,
		Destination: Destination,
	}
	service.Find(cfg)
}

func parseArgs() {
	findCmd.Flags().StringVarP(&Target, "Target", "t", "the target path", "the target path")
	findCmd.Flags().StringVarP(&Destination, "Destination", "d", "the Destination path", "the Destination path")
	findCmd.Flags().BoolVarP(&Verbose, "Verbose", "v", false, "visible the process")
}
