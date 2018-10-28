package controllers

import (
	"github.com/go-pg/pg"
	"github.com/labstack/echo"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
	"regexp"
)

type SignupStruct struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func UserSignup(c echo.Context) error {
	db := c.Get("DB").(*pg.DB)
	signup := new(SignupStruct)
	if err := c.Bind(signup); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9_-]{1,39}$`)
	if !regex.MatchString(signup.Username) {
		err := models.JsonResponse{
			Code: http.StatusUnprocessableEntity,
			Err:  "username may only contain alphanumeric, - and _ characters and be less than 40 caracters",
		}
		return c.JSON(err.Code, err)
	}
	user := new(models.User)
	if err := db.Model(user).Column("id").Where("username = ?", signup.Username).First(); err == nil {
		c.Logger().Error(err)
		err := models.JsonResponse{
			Code: http.StatusConflict,
			Err:  "username is already taken",
		}
		return c.JSON(err.Code, err)
	}
	user, err := dao.SignupUser(signup.Username, signup.Email, signup.Password)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	jwt, err := dao.GenerateToken(user)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": jwt,
	})
}

func UserLogin(c echo.Context) error {
	login := new(LoginStruct)
	if err := c.Bind(login); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusUnprocessableEntity, nil)
	}
	user := new(models.User)
	if err := dao.GetUserByUsername(user, login.Username); err != nil {
		c.Logger().Error(err)
		err := models.JsonResponse{
			Code: http.StatusNotFound,
			Err:  "user not found",
		}
		return c.JSON(err.Code, err)
	}
	if err := user.VerifyPassword(login.Password); err != nil {
		c.Logger().Error(err)
		err := models.JsonResponse{
			Code: http.StatusUnauthorized,
			Err:  "password invalid",
		}
		return c.JSON(err.Code, err)
	}
	jwt, err := dao.GenerateToken(user)
	if err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, map[string]string{
		"token": jwt,
	})
}

func GetUserProfile(c echo.Context) error {
	username := c.Param("username")
	user := new(models.User)
	err := dao.GetProfileByUsername(user, username)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, user)
}

func GetCurrentUserProfile(c echo.Context) error {
	current_user := c.Get("current_user").(*models.User)
	err := dao.GetProfile(current_user)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, current_user)
}

func GetCurrentUserFriends(c echo.Context) error {
	current_user := c.Get("current_user").(*models.User)
	err := dao.GetUserWithFriends(current_user)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, current_user.Friends)
}

func AddFriend(c echo.Context) error {
	current_user := c.Get("current_user").(*models.User)
	username := c.Param("username")
	friend := new(models.User)
	if err := dao.GetUserByUsername(friend, username); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	if err := dao.AddFriend(current_user, friend); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func RemoveFriend(c echo.Context) error {
	current_user := c.Get("current_user").(*models.User)
	username := c.Param("username")
	friend := new(models.User)
	if err := dao.GetUserByUsername(friend, username); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusNotFound)
	}
	if err := dao.RemoveFriend(current_user, friend); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
