package cmd

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tendresse",
	Short: "Tendresse, best app EU",
	Long:  `Tendresse, gif sharing platform, best app EU`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
