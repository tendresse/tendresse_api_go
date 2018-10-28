package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tendresse/tendresse_api_go/app"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
)

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testLoadCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Cli for testing management",
	Long:  `Load some test data`,
}

var testLoadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load test data",
	Run: func(cmd *cobra.Command, args []string) {
		app.InitEnv()
		database.Init()
		defer database.CloseDB()
		database.ResetDB()
		user1, _ := dao.SignupUser("user1", "user1@mail.com", "user1")
		user2, _ := dao.SignupUser("user2", "user2@mail.com", "user2")
		user3, _ := dao.SignupUser("user3", "user3@mail.com", "user3")
		dao.AddFriend(user1, user2)
		gif1 := new(models.Gif)
		gif1.Url = "https://66.media.tumblr.com/8491150b6c9bc08cf32c4068e4f25b5f/tumblr_pgsrxx4d6u1vgwlcko1_1280.gif"
		dao.AddGif(gif1)
		dao.UpdateGifTags(gif1, []string{"tag1", "tag2"})
		achievement1 := &models.Achievement{
			Name:      "achievement1",
			Condition: 1,
			Xp:        10,
			Type:      "receive",
			Icon:      "img/WW04imNRnDkUJx5DkANA_fist.png",
			TagID:     1,
		}
		achievement2 := &models.Achievement{
			Name:      "achievement2",
			Condition: 1,
			Xp:        10,
			Type:      "receive",
			Icon:      "img/WW04imNRnDkUJx5DkANA_fist.png",
			TagID:     2,
		}
		dao.AddAchievement(achievement1)
		dao.AddAchievement(achievement2)
		dao.SendTendresse(user2, user1)
		dao.SendTendresse(user3, user1)
		admin_role := new(models.Role)
		admin_role.Name = "admin"
		dao.GetOrCreateRole(admin_role)
		dao.AddRoleToUser(admin_role, user1)
	},
}
