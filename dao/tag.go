package dao

import (
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
)

func GetOrCreateTag(tag *models.Tag) error {
	db := database.GetDB()
	banned_tags := []string{
		".gifs", "erotic", "famous", "fuck", "fucking", "funny", "more", "info", "animated", "have", "week", "gif", "gifs", "horny", "hot",
		"huge", "naughty", "nsfw", "porn", "porn gif", "porno", "pornstar", "pr0n", "scene", "sex", "sexe",
		"sexy", "star", "the", "tumblr", "wild", "xxx",
	}
	for _, banned_tag := range banned_tags {
		if tag.Name == banned_tag {
			return errors.New("this tag is banned")
		}
	}
	_, err := db.Model(tag).
		Column("id").
		Where("name = ?", tag.Name).
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("id").
		SelectOrInsert()
	return err
}

func GetOrCreateTags(s_tags []string) []*models.Tag {
	var tags []*models.Tag
	for _, s_tag := range s_tags {
		tag := new(models.Tag)
		tag.Name = s_tag
		if err := GetOrCreateTag(tag); err != nil {
			log.Debug(err)
			continue
		}
		tags = append(tags, tag)
	}
	return tags
}

//
//func GetBannedTags(db *pg.DB, tags []*models.Tag) error {
//	count, err := db.Model(&models.Tag{}).Count()
//	if err != nil {
//		return err
//	}
//	return db.Model(tags).
//		Where("banned = ?", true).
//		Limit(count).
//		Select()
//}
//
//func GetTag(db *pg.DB, tag *models.Tag) error {
//	return db.Select(&tag)
//}
//func GetTags(db *pg.DB, tags []*models.Tag) error {
//	return db.Model(&tags).Select()
//}
//
//func GetAllTags(db *pg.DB, tags *[]models.Tag) error {
//	count, err := db.Model(&models.Tag{}).Count()
//	if err != nil {
//		return err
//	}
//	return db.Model(&tags).Limit(count).Select()
//}
//
//func UpdateTag(db *pg.DB, tag *models.Tag) error {
//	return db.Update(tag)
//}
//
//func DeleteTag(db *pg.DB, tag *models.Tag) error {
//	return db.Delete(&tag)
//}
//func DeleteTags(db *pg.DB, tags []*models.Tag) error {
//	for _, tag := range tags {
//		if err := c.DeleteTag(tag); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func GetFullTag(db *pg.DB, tag *models.Tag) error {
//	return db.Model(&tag).Column("tag.*", "Gifs").Column("tag.*", "Achievements").Where("tag.id = ?", tag.Id).First()
//}
//func GetFullTags(db *pg.DB, tags []*models.Tag) error {
//	return db.Model(&tags).Column("tag.*", "Gifs").Column("tag.*", "Achievements").Select()
//}
//
//func GetTagByTitle(db *pg.DB, title string, tag *models.Tag) error {
//	return db.Model(&tag).Where("title = ?", title).Select()
//}
//
func AddTagToGif(tag *models.Tag, gif *models.Gif) error {
	db := database.GetDB()
	gifs_tags := new(models.GifsTags)
	gifs_tags.GifID = gif.ID
	gifs_tags.TagID = tag.ID
	err := db.Insert(gifs_tags)
	return errors.Wrap(err, "adding tag to gif")
}

func RemoveTagFromGif(tag *models.Tag, gif *models.Gif) error {
	db := database.GetDB()
	gt := new(models.GifsTags)
	gt.TagID = tag.ID
	gt.GifID = gif.ID
	err := db.Delete(gt)
	return errors.Wrap(err, "remove tag from gif")
}

func GetTags(tags *[]*models.Tag) error {
	db := database.GetDB()
	err := db.Model(tags).Select()
	return errors.Wrap(err, "get all tags")
}

func DeleteTag(tag *models.Tag) error {
	db := database.GetDB()
	err := db.Delete(tag)
	return errors.Wrap(err, "delete tag")
}

func GetFullTag(tag *models.Tag) error {
	db := database.GetDB()
	err := db.Model(tag).Where("tag.id = ?", tag.ID).
		Column("tag.*", "Gifs").
		Column("tag.*", "Achievements").
		First()
	return errors.Wrap(err, "get full tag : tag + gifs")
}
