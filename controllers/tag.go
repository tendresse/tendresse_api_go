package controllers

import (
	"github.com/labstack/echo"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
	"strconv"
)

func GetTag(c echo.Context) error {
	tag_id, err := strconv.Atoi(c.Param("tag_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	tag := new(models.Tag)
	tag.ID = tag_id
	if err := dao.GetFullTag(tag); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, tag)
}

func GetTags(c echo.Context) error {
	tags := new([]*models.Tag)
	if err := dao.GetTags(tags); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, tags)
}

func DeleteTag(c echo.Context) error {
	tag_id, err := strconv.Atoi(c.Param("tag_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	tag := new(models.Tag)
	tag.ID = tag_id
	if err := dao.DeleteTag(tag); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
