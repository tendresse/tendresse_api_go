package controllers

import (
	"github.com/labstack/echo"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
	"net/url"
	"strconv"
)

type newBlog struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

func GetBlogs(c echo.Context) error {
	blogs := new([]*models.Blog)
	if err := dao.GetBlogs(blogs); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, blogs)
}

func AddBlog(c echo.Context) error {
	new_blog := new(newBlog)
	if err := c.Bind(new_blog); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	full_url, err := url.Parse(new_blog.Url)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	if full_url.Host == "" {
		c.Logger().Error("empty URL for blog")
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	blog := new(models.Blog)
	blog.Url = full_url.Host
	if err := dao.GetBlogByUrl(blog, blog.Url); err == nil {
		c.Logger().Debug(err)
		return c.NoContent(http.StatusConflict)
	}
	blog.Description = new_blog.Description
	if err := dao.AddBlog(blog); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	go dao.ImportGifsFromBlog(blog)
	return c.JSON(http.StatusOK, blog)
}

func DeleteBlog(c echo.Context) error {
	blog_id, err := strconv.Atoi(c.Param("blog_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	blog := new(models.Blog)
	blog.ID = blog_id
	if err := dao.DeleteBlog(blog); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func DeleteBlogAndGifs(c echo.Context) error {
	blog_id, err := strconv.Atoi(c.Param("blog_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	blog := new(models.Blog)
	blog.ID = blog_id
	if err := dao.DeleteGifsByBlog(blog); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if err := dao.DeleteBlog(blog); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
