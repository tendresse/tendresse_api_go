package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tendresse/tendresse_api_go/app"
	"github.com/tendresse/tendresse_api_go/database"
	"log"
	"sort"
)

func init() {
	dbCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateListCmd)
	migrateCmd.AddCommand(migrateApplyCmd)
	migrateCmd.AddCommand(migrateRollbackCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage DB migrations",
}

var migrateApplyCmd = &cobra.Command{
	Use:   "apply [migration_id]",
	Short: "Apply a migration by ID",
	Long:  `Launch a migration, for example "db migrate apply 3 4 5"`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitEnv()
		database.Init()
		db := database.GetDB()
		defer database.CloseDB()
		migrations := getMigrations()
		for _, arg := range args {
			if mig, present := migrations[arg]; present {
				mig.migrate(db)
			} else {
				log.Println("error, 404 migration not found")
				log.Println("please use 'list' to show the list of all migrations")
			}
		}
	},
}

var migrateRollbackCmd = &cobra.Command{
	Use:   "rollback [migrations_id]",
	Short: "Rollback migrations by ID",
	Long:  `Rollback migrations, for example "db migrate rollback 1 2 3"`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitEnv()
		database.Init()
		db := database.GetDB()
		defer database.CloseDB()
		migrations := getMigrations()
		for _, arg := range args {
			if mig, present := migrations[arg]; present {
				mig.rollback(db)
			} else {
				log.Println("error, 404 migration not found")
				log.Println("please use 'list' to show the list of all migrations")
			}
		}
	},
}

var migrateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all migrations",
	Run: func(cmd *cobra.Command, args []string) {
		ms := getMigrations()
		if len(ms) < 1 {
			log.Println("no migrations")
			return
		}
		// TODO refactor ce sorting de map degeulasse
		migrations := make([]string, 0, len(ms))
		for mig := range ms {
			migrations = append(migrations, mig)
		}
		sort.Strings(migrations)
		for mig := range migrations {
			mk := migrations[mig]
			mv := ms[mk]
			log.Println(mk, "|", mv.description)
		}
	},
}
