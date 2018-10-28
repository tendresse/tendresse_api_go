package controllers

import (
	"github.com/labstack/echo"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
	"strconv"
)

func GetRandomGif(c echo.Context) error {
	gif := new(models.Gif)
	if err := dao.GetRandomGif(gif); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, gif)
}

func DeleteGif(c echo.Context) error {
	gif_id, err := strconv.Atoi(c.Param("gif_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	gif := new(models.Gif)
	gif.ID = gif_id
	if err := dao.DeleteGif(gif); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

type gifJson struct {
	Url string `json:"url"`
	tagsJson
}

type tagsJson struct {
	Tags []string `json:"tags"`
}

func AddGif(c echo.Context) error {
	gif_json := new(gifJson)
	if err := c.Bind(gif_json); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	gif := new(models.Gif)
	if err := dao.GetGifByUrl(gif, gif_json.Url); err == nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusConflict)
	}
	gif.Url = gif_json.Url
	if err := dao.AddGif(gif); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	dao.UpdateGifTags(gif, gif_json.Tags)
	return c.JSON(http.StatusOK, gif)
}

func UpdateGifTags(c echo.Context) error {
	gif_id, err := strconv.Atoi(c.Param("gif_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	gif := new(models.Gif)
	gif.ID = gif_id
	if err := dao.GetFullGif(gif); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	new_tags := new(tagsJson)
	if err := c.Bind(new_tags); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	dao.UpdateGifTags(gif, new_tags.Tags)
	return c.NoContent(http.StatusOK)
}

func GetGifs(c echo.Context) error {
	gifs := new([]*models.Gif)
	if err := dao.GetGifs(gifs); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, gifs)
}

func CheckAllGifs(c echo.Context) error {
	err := dao.CheckAllGifs()
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
