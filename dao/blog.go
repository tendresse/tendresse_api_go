package dao

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//import (
//	"github.com/go-pg/pg"
//	"github.com/tendresse/tendresse_api_go/models"
//)
//
//func CreateBlog(db *pg.DB, blog *models.Blog) error {
//	return db.Insert(blog)
//}
//func CreateBlogs(db *pg.DB, blogs []*models.Blog) error {
//	for _, blog := range blogs {
//		if err := c.CreateBlog(blog); err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func UpdateBlog(db *pg.DB, blog *models.Blog) error {
//	return db.Update(blog)
//}
//

func GetBlogByUrl(blog *models.Blog, url string) error {
	db := database.GetDB()
	err := db.Model(blog).Where(`blog.url = ?`, url).First()
	return errors.Wrap(err, "get blog by url")
}

func AddBlog(blog *models.Blog) error {
	db := database.GetDB()
	err := db.Insert(blog)
	return errors.Wrap(err, "inserting blog")
}

func fetchTumblr(url string, tumblr *models.Tumblr) error {
	http_client := &http.Client{Timeout: 10 * time.Second}
	resp, err := http_client.Get(url)
	defer resp.Body.Close()
	if err != nil {
		err = errors.Wrap(err, "cannot fetch tumblr API")
		log.Error(err)
		return err
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "cannot read body of tumblr API")
		log.Error(err)
		return err
	}

	err = json.Unmarshal(contents, &tumblr)
	if err != nil {
		err = errors.Wrap(err, "cannot unmarshal json response of tumblr API")
		log.Error(err)
		return err
	}

	if tumblr.Meta.Status != 200 {
		err = errors.Wrap(err, "this tumblr is not valid")
		log.Error(err)
		return err
	}
	return nil
}

func ImportGifsFromBlog(blog *models.Blog) {
	api_blog_url := strings.Join([]string{
		"https://api.tumblr.com/v2/blog/",
		blog.Url,
		"/posts/photo?api_key=",
		os.Getenv("TUMBLR_API_KEY"),
	}, "")
	tumblr := new(models.Tumblr)
	if err := fetchTumblr(api_blog_url, tumblr); err != nil {
		log.Error(err)
		return
	}
	api_blog_url = strings.Join([]string{api_blog_url, "&offset="}, "")

	for i := 0; i < tumblr.Response.Blog.TotalPosts; i += 20 {
		posts_url := strings.Join([]string{api_blog_url, strconv.Itoa(i)}, "")
		tumblr = new(models.Tumblr)

		if err := fetchTumblr(posts_url, tumblr); err != nil {
			log.Error(err)
			return
		}

		for _, post := range tumblr.Response.Posts {
			if post.Type != "photo" {
				continue
			}
			gif_url := post.Photos[0].OriginalSize.Url
			if gif_url[len(gif_url)-3:] != "gif" {
				continue
			}
			gif := new(models.Gif)
			gif.Url = gif_url
			gif.BlogID = blog.ID
			if err := AddGif(gif); err != nil {
				log.Error(err)
				continue
			}
		Tags:
			for _, ltags := range post.Tags {
				for _, tag_name := range strings.Fields(ltags) {
					if len(tag_name) <= 3 {
						continue Tags
					}
					tag := new(models.Tag)
					tag.Name = tag_name
					if err := GetOrCreateTag(tag); err != nil {
						log.Error(err)
						continue Tags
					}
					if err := AddTagToGif(tag, gif); err != nil {
						log.Error(err)
					}
				}
			}
		}
	}
}

func GetBlogs(blogs *[]*models.Blog) error {
	db := database.GetDB()
	err := db.Model(blogs).Select()
	return errors.Wrap(err, "get all blogs")
}

//
//func GetAllBlogs(db *pg.DB, blogs *[]models.Blog) error {
//	count, err := db.Model(&models.Blog{}).Count()
//	if err != nil {
//		return err
//	}
//	return db.Model(&blogs).Limit(count).Select()
//}
//
func GetFullBlog(blog *models.Blog) error {
	db := database.GetDB()
	err := db.Model(blog).Where("id = ?", blog.ID).
		Column("blog.*", "Gifs").
		First()
	return errors.Wrap(err, "get full blog: blog + gifs")
}

//func GetFullBlogs(db *pg.DB, blogs []*models.Blog) error {
//	return db.Model(&blogs).Column("blog.*", "Gifs").Select()
//}
//
//func GetBlogByTitle(db *pg.DB, title string, blog *models.Blog) error {
//	return db.Model(&blog).Where("title = ?", title).First()
//}
//
//func GetBlogByUrl(db *pg.DB, url string, blog *models.Blog) error {
//	return db.Model(&blog).Where("url = ?", url).First()
//}
//
func DeleteBlog(blog *models.Blog) error {
	db := database.GetDB()
	err := db.Delete(blog)
	return errors.Wrap(err, "delete blog by string ID")
}

//func DeleteBlogs(db *pg.DB, blogs []*models.Blog) error {
//	// TODO : delete cascade on blogs delete
//	for _, blog := range blogs {
//		if err := c.DeleteBlog(blog); err != nil {
//			return err
//		}
//	}
//	return nil
//}
