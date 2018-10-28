package controllers

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
	"strconv"
)

func SendTendresse(c echo.Context) error {
	current_user := c.Get("current_user").(*models.User)
	username := c.Param("username")
	friend := new(models.User)
	if err := dao.GetUserByUsername(friend, username); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	if err := dao.SendTendresse(current_user, friend); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func GetTendresses(c echo.Context) error {
	current_user := c.Get("current_user").(*models.User)
	tendresses := new([]*models.Tendresse)
	err := dao.GetPendingTendresses(tendresses, current_user)
	if err != nil {
		c.Logger().Error(errors.Wrap(err, "get pending tendresses of current_user"))
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, tendresses)
}

func ChangeTendresseAsViewed(c echo.Context) error {
	current_user := c.Get("current_user").(*models.User)
	tendresse_id, err := strconv.Atoi(c.Param("tendresse_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	tendresse := new(models.Tendresse)
	tendresse.ID = tendresse_id
	if err := dao.GetTendresse(tendresse); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	if tendresse.ReceiverID != current_user.ID {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"err": "this tendresse was not sent to you",
		})
	}
	if err := dao.ChangeTendresseAsViewed(tendresse); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"err": "can't change state of tendresse as viewed",
		})
	}
	return c.NoContent(http.StatusOK)
}
