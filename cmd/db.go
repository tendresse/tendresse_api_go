package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tendresse/tendresse_api_go/app"
	"github.com/tendresse/tendresse_api_go/database"
)

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(dbResetCmd)
	dbCmd.AddCommand(dbDropCmd)
	dbCmd.AddCommand(dbCreateCmd)
}

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Cli for DB management",
	Long:  `Manage the DB like init and migrate`,
}

var dbResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all tables",
	Long:  `Drop && Auto-Migrate models`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitEnv()
		database.Init()
		defer database.CloseDB()
		database.ResetDB()
	},
}

var dbDropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Drop all tables",
	Run: func(cmd *cobra.Command, args []string) {
		app.InitEnv()
		database.Init()
		defer database.CloseDB()
		database.DropDB()
	},
}

var dbCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create all tables",
	Run: func(cmd *cobra.Command, args []string) {
		app.InitEnv()
		database.Init()
		defer database.CloseDB()
		database.CreateDB()
	},
}
