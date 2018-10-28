package cmd

import (
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/tendresse/tendresse_api_go/app"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
	"regexp"
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(createAdminCmd)
	createAdminCmd.Flags().StringVarP(&username, "username", "u", "", "username of the admin")
	createAdminCmd.Flags().StringVarP(&email, "email", "e", "", "email of the admin")
	createAdminCmd.Flags().StringVarP(&password, "password", "p", "", "password of the admin")
	createAdminCmd.MarkFlagRequired("username")
	createAdminCmd.MarkFlagRequired("email")
	createAdminCmd.MarkFlagRequired("password")
}

var username string
var email string
var password string

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Cli for User management",
	Long:  `Manage the Users`,
}

var createAdminCmd = &cobra.Command{
	Use:   "admin-create",
	Short: "Reset all tables",
	Long:  `Drop && Auto-Migrate models`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitEnv()
		database.Init()
		defer database.CloseDB()
		regex := regexp.MustCompile(`^[a-zA-Z0-9_-]{1,39}$`)
		if !regex.MatchString(username) {
			log.Error("username may only contain alphanumeric, - and _ characters and be less than 40 caracters")
			return
		}
		user, err := dao.SignupUser(username, email, password)
		if err != nil {
			log.Error(err)
			return
		}
		admin_role := new(models.Role)
		admin_role.Name = "admin"
		if err := dao.GetOrCreateRole(admin_role); err != nil {
			log.Error(err)
			return
		}
		if err := dao.AddRoleToUser(admin_role, user); err != nil {
			log.Error(err)
			return
		}
		log.Printf("admin user %q was created.", user.Username)
	},
}
