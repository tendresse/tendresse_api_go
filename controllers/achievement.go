package controllers

import (
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/dao"
	"github.com/tendresse/tendresse_api_go/models"
	"net/http"
	"strconv"
)

type achievementJson struct {
	Name      string `json:"name"`
	Condition int    `json:"condition"`
	Icon      string `json:"icon"`
	Type      string `json:"type"`
	Xp        int    `json:"xp"`
	Tag       string `json:"tag"`
}

func (aJson achievementJson) Validate() error {
	if aJson.Name == "" {
		return errors.New("Name cannot be empty")
	}
	if aJson.Condition < 1 {
		return errors.New("Condition cannot be < 1")
	}
	if aJson.Type != "send" && aJson.Type != "receive" {
		return errors.New("Type must be 'receive' or 'send'")
	}
	if aJson.Icon == "" {
		return errors.New("Icon cannot be empty")
	}
	if aJson.Xp < 1 {
		return errors.New("Xp cannot be < 1")
	}
	if aJson.Tag == "" {
		return errors.New("Tag cannot be empty")
	}
	return nil
}

func AddAchievement(c echo.Context) error {
	aJson := new(achievementJson)
	if err := c.Bind(aJson); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	achievement := new(models.Achievement)
	if err := dao.GetAchievementByName(achievement, aJson.Name); err == nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusConflict)
	}
	if err := aJson.Validate(); err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	achievement.Name = aJson.Name
	achievement.Condition = aJson.Condition
	achievement.Xp = aJson.Xp
	achievement.Type = aJson.Type
	achievement.Icon = aJson.Icon
	tag := new(models.Tag)
	tag.Name = aJson.Tag
	// check if this achievement is possible
	// aka if there is N(condition) different Gif with this Tag
	if err := dao.GetOrCreateTag(tag); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	achievement.TagID = tag.ID
	achievement.Tag = tag
	if err := dao.AddAchievement(achievement); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	IsAchievementPossible(achievement)
	return c.JSON(http.StatusOK, achievement)
}

func IsAchievementPossible(achievement *models.Achievement) error {
	count, err := dao.CountGifsWithTag(achievement.Tag)
	if err != nil {
		return err
	}
	if count < achievement.Condition {
		return errors.New("this achievement is not possible with distinct Gifs")
	}
	return nil
}

func GetAchievements(c echo.Context) error {
	achievements := new([]*models.Achievement)
	if err := dao.GetAchievements(achievements); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, achievements)
}

func GetAchievementsLimited(c echo.Context) error {
	achievements := new([]*models.Achievement)
	if err := dao.GetAchievementsLimited(achievements); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, achievements)
}

func DeleteAchievement(c echo.Context) error {
	achievement_id, err := strconv.Atoi(c.Param("achievement_id"))
	if err != nil {
		return c.NoContent(http.StatusUnprocessableEntity)
	}
	achievement := new(models.Achievement)
	achievement.ID = achievement_id
	if err := dao.DeleteAchievement(achievement); err != nil {
		c.Logger().Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
