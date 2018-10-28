package dao

import (
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
)

func SendTendresse(sender *models.User, receiver *models.User) error {
	db := database.GetDB()
	gif := new(models.Gif)
	if err := GetRandomGif(gif); err != nil {
		return err
	}
	tendresse := new(models.Tendresse)
	tendresse.SenderID = sender.ID
	tendresse.ReceiverID = receiver.ID
	tendresse.GifID = gif.ID
	if err := db.Insert(tendresse); err != nil {
		return errors.Wrap(err, "inserting tendresse after creation")
	}
	achievements := []*models.Achievement{}
	err := GetAchievementsFromGif(&achievements, gif)
	if err != nil {
		return err
	}
	if len(achievements) > 0 {
		go UpdateUserAchievements(receiver, achievements, "receive")
		sender_achievements := []*models.Achievement{}
		err := GetSenderAchievementsWithoutTag(&sender_achievements)
		if err != nil {
			log.Error(err)
		}
		achievements = append(achievements, sender_achievements...)
		go UpdateUserAchievements(sender, achievements, "send")
	}
	return nil
}

func GetTendresse(tendresse *models.Tendresse) error {
	db := database.GetDB()
	err := db.Model(tendresse).
		Where(`tendresse.id = ?`, tendresse.ID).
		First()
	return errors.Wrap(err, "get tendresse")
}

func GetPendingTendresses(tendresses *[]*models.Tendresse, user *models.User) error {
	db := database.GetDB()
	err := db.Model(tendresses).
		Where("tendresse.receiver_id = ?", user.ID).
		Where("tendresse.viewed IS NOT TRUE").
		// Column("tendresse.*", "Receiver").
		Column("tendresse.*", "Sender").
		Column("tendresse.*", "Gif").
		Select()
	return errors.Wrap(err, "selecting tendresses not viewed")

}

func ChangeTendresseAsViewed(tendresse *models.Tendresse) error {
	db := database.GetDB()
	tendresse.Viewed = true
	err := db.Update(tendresse)
	return errors.Wrap(err, "updating tendresse to change viewed to True")
}

//
//func CountSenderTendresses(db *pg.DB, sender_id int) (int, error) {
//	return db.Model(&models.Tendresse{}).Where("sender_id = ?", sender_id).Count()
//}
