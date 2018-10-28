package dao

import (
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
	"time"
)

func AddGif(gif *models.Gif) error {
	db := database.GetDB()
	err := db.Insert(gif)
	return errors.Wrap(err, "inserting gif")
}

func GetRandomGif(gif *models.Gif) error {
	db := database.GetDB()
	// get random tag
	tag_id := 0
	gif_id := 0
	// next query fails if number of gifs is very low
	_, err := db.QueryOne(&tag_id, "SELECT t.id FROM tags t OFFSET FLOOR(RANDOM() * (SELECT COUNT(*) FROM tags)) LIMIT 1 ;")
	if err != nil {
		return errors.Wrap(err, "random select of a Tag")
	}
	// get random gif from random tag
	_, err = db.QueryOne(&gif_id, "SELECT gif_id FROM gifs_tags WHERE tag_id = ? OFFSET FLOOR( RANDOM() * (SELECT COUNT(*) FROM gifs_tags WHERE tag_id = ?)) LIMIT 1 ;", tag_id, tag_id)
	if err != nil {
		return errors.Wrap(err, "random select of a Gif after random select of a Tag")
	}
	gif.ID = gif_id
	return GetFullGif(gif)
}

func GetGifs(gifs *[]*models.Gif) error {
	db := database.GetDB()
	err := db.Model(gifs).Select()
	return errors.Wrap(err, "get gifs")
}

//func GetAllGifs(db *pg.DB, gifs *[]models.Gif) error {
//	count, err := db.Model(&models.Gif{}).Count()
//	if err != nil {
//		return err
//	}
//	return db.Model(&gifs).Limit(count).Select()
//}
//
func DeleteGif(gif *models.Gif) error {
	db := database.GetDB()
	err := db.Delete(gif)
	return errors.Wrap(err, "deleting gif")
}

func GetFullGif(gif *models.Gif) error {
	db := database.GetDB()
	err := db.Model(gif).Where("gif.id = ?", gif.ID).
		Column("gif.*", "Blog").
		Column("gif.*", "Tags").
		First()
	return errors.Wrap(err, "get full gif : tags + blog")
}

func GetGifByUrl(gif *models.Gif, url string) error {
	db := database.GetDB()
	err := db.Model(gif).
		Column("gif.id").
		Where("gif.url = ?", url).
		First()
	return errors.Wrap(err, "getting gif by url")
}

func UpdateGifTags(gif *models.Gif, new_tags []string) {
CurrentTags:
	for _, tag := range gif.Tags {
		for _, new_tag := range new_tags {
			if new_tag == tag.Name {
				continue CurrentTags
			}
		}
		if err := RemoveTagFromGif(&tag, gif); err != nil {
			log.Error(err)
		}
	}
NewTags:
	for _, new_tag := range new_tags {
		for _, tag := range gif.Tags {
			if new_tag == tag.Name {
				continue NewTags
			}
		}
		tag := new(models.Tag)
		tag.Name = new_tag
		if err := GetOrCreateTag(tag); err != nil {
			log.Error(err)
			continue NewTags
		}
		if err := AddTagToGif(tag, gif); err != nil {
			log.Error(err)
			continue NewTags
		}
	}
}

func CountGifsWithTag(tag *models.Tag) (int, error) {
	db := database.GetDB()
	count, err := db.Model((*models.Gif)(nil)).
		Column("gif.id").
		Join("JOIN gifs_tags AS g ON g.gif_id = gif.id").
		Where("g.tag_id = ?", tag.ID).
		Count()
	return count, errors.Wrap(err, "counting gifs that have a certain tag")
}

func DeleteGifsByBlog(blog *models.Blog) error {
	db := database.GetDB()
	_, err := db.Model((*models.Gif)(nil)).
		Where("gif.blog_id = ?", blog.ID).
		Delete()
	return errors.Wrap(err, "delete gifs by string blog ID")
}

func CheckAllGifs() error {
	db := database.GetDB()
	if err := CheckTumblr(); err != nil {
		return err
	}
	gifs := new([]*models.Gif)
	err := db.Model(gifs).Select()
	if err != nil {
		return err
	}
	for _, gif := range *gifs {
		go CheckGif(gif)
	}
	return nil
}

func CheckGif(gif *models.Gif) {
	http_client := &http.Client{Timeout: 10 * time.Second}
	resp, err := http_client.Head(gif.Url)
	defer resp.Body.Close()
	if err != nil {
		log.Error(err)
	}
	if resp.StatusCode >= 300 {
		log.Errorf("gif %d with url %q is not available - deleting", gif.ID, gif.Url)
		DeleteGif(gif)
	}
}

func CheckTumblr() error {
	http_client := &http.Client{Timeout: 10 * time.Second}
	resp, err := http_client.Head("https://www.tumblr.com/")
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return errors.New("tumblr seems down")
	}
	return nil
}
