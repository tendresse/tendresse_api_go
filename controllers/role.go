package controllers

import (
	"github.com/labstack/echo"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
)

func GetRoles(c echo.Context) error {
	roles := new([]*models.Role)
	if err := dao.GetRoles(roles); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, roles)
}

func AddRoleToUser(c echo.Context) error {
	username := c.Param("username")
	user := new(models.User)
	if err := dao.GetUserByUsername(user, username); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	role_name := c.Param("role_name")
	for _, user_role := range user.Roles {
		if user_role.Name == role_name {
			return c.NoContent(http.StatusOK)
		}
	}
	role := new(models.Role)
	role.Name = role_name
	if err := dao.GetOrCreateRole(role); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	if err := dao.AddRoleToUser(role, user); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func RemoveRoleFromUser(c echo.Context) error {
	username := c.Param("username")
	user := new(models.User)
	if err := dao.GetUserByUsername(user, username); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	role_name := c.Param("role_name")
	for _, user_role := range user.Roles {
		if user_role.Name == role_name {
			role := new(models.Role)
			role.Name = role_name
			if err := dao.GetOrCreateRole(role); err != nil {
				c.Logger().Error(err)
				return c.NoContent(http.StatusNotFound)
			}
			if err := dao.RemoveRoleFromUser(role, user); err != nil {
				c.Logger().Error(err)
				return c.NoContent(http.StatusInternalServerError)
			}
			return c.NoContent(http.StatusOK)
		}
	}
	return c.NoContent(http.StatusOK)
}
