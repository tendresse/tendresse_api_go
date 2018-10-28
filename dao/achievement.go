package dao

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
)

//
//import (
//	"github.com/go-pg/pg"
//	"github.com/tendresse/tendresse_api_go/models"
//)

func GetAchievementsFromGif(achievements *[]*models.Achievement, gif *models.Gif) error {
	/*
	 * trouver les achievements liés au Gif de façon efficace
	 * Gif -> Tags -> Achievements
	 * idée : faire un array des ID des tags et utiliser le SELECT "in" de go-pg
	 */
	db := database.GetDB()
	var id_tags []int
	for _, tag := range gif.Tags {
		id_tags = append(id_tags, tag.ID)
	}
	err := db.Model(achievements).
		Where("tag_id in (?)", pg.In(id_tags)).
		Select()
	return errors.Wrap(err, "select achievements with pg.In tags_ids")
}

func AddAchievement(achievement *models.Achievement) error {
	db := database.GetDB()
	err := db.Insert(achievement)
	return errors.Wrap(err, "add achievement")
}

func GetAchievements(achievements *[]*models.Achievement) error {
	db := database.GetDB()
	err := db.Model(achievements).Select()
	return errors.Wrap(err, "get all achievements")
}

func GetAchievementsLimited(achievements *[]*models.Achievement) error {
	db := database.GetDB()
	err := db.Model(achievements).
		Column("id").
		Column("name").
		Select()
	return errors.Wrap(err, "get all achievements limited")
}

func DeleteAchievement(achievement *models.Achievement) error {
	db := database.GetDB()
	err := db.Delete(achievement)
	return errors.Wrap(err, "delete achievement")
}

//
//func CreateAchievements(db *pg.DB, achievements []*models.Achievement) error {
//	for _, achievement := range achievements {
//		if err := CreateAchievement(db, achievement); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func UpdateAchievement(db *pg.DB, achievement *models.Achievement) error {
//	return db.Update(achievement)
//}
//
//func GetAchievement(db *pg.DB, achievement *models.Achievement) error {
//	return db.Select(&achievement)
//}
//func GetAchievements(db *pg.DB, achievements []*models.Achievement) error {
//	return db.Model(&achievements).Select()
//}
//
//func GetAllAchievements(db *pg.DB, achievements *[]models.Achievement) error {
//	count, err := db.Model(&models.Achievement{}).Count()
//	if err != nil {
//		return err
//	}
//	return db.Model(&achievements).Limit(count).Select()
//}
//
//func GetFullAchievement(db *pg.DB, achievement *models.Achievement) error {
//	return db.Model(&achievement).Column("achievement.*", "Tag").Where("id = ?", achievement.Id).First()
//}
//func GetFullAchievements(achievements []*models.Achievement) error {
//	return db.Model(&achievements).Column("achievement.*", "Tag").Select()
//}
//
func GetAchievementByName(achievement *models.Achievement, name string) error {
	db := database.GetDB()
	err := db.Model(achievement).Column("id").Where("name = ?", name).First()
	return errors.Wrap(err, "get achievement by name")
}

//
//func GetOrCreateAchievement(db *pg.DB, achievement *models.Achievement) error {
//	return db.Select(&achievement)
//}
//
//func DeleteAchievement(db *pg.DB, achievement *models.Achievement) error {
//	return db.Delete(&achievement)
//}
//func DeleteAchievements(db *pg.DB, achievements []*models.Achievement) error {
//	for _, achievement := range achievements {
//		if err := c.DeleteAchievement(achievement); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
func GetSenderAchievementsWithoutTag(achievements *[]*models.Achievement) error {
	db := database.GetDB()
	err := db.Model(achievements).
		Where("achievement.type = ?", "sender").
		Where("achievement.tag_id IS null").
		Select()
	return errors.Wrap(err, "getting sender achievements without tag from DB")
}
