package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tendresse/tendresse_api_go/app"
)

func init() {
	rootCmd.AddCommand(webCmd)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "launch the web server",
	Long:  `Launch the echo web-server`,
	Run: func(cmd *cobra.Command, args []string) {
		app.Launch()
	},
}
