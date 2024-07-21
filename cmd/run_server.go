package cmd

import (
	"github.com/spf13/cobra"

	"core/configs"
	"core/jobs"
	"core/server"
	"core/third_party"
)

var cmdRunServer = &cobra.Command{
	Use:   "run_server",
	Short: "Runs the gin api server",
	Long:  `Runs the gin api server`,
	Run: func(cmd *cobra.Command, args []string) {
		configs.Env{}.Initialize()
		third_party.Supertoken{}.Initialize()
		jobs.Initialize()
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(cmdRunServer)
}
