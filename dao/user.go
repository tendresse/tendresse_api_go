package dao

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/tendresse/tendresse_api_go/database"
	"github.com/tendresse/tendresse_api_go/models"
)

func SignupUser(username string, email string, password string) (*models.User, error) {
	db := database.GetDB()
	new_user := &models.User{
		Username: username,
		Email:    email,
	}
	if err := new_user.SetPassword(password); err != nil {
		return nil, errors.Wrap(err, "encrypting signup user password")
	}
	if err := db.Insert(new_user); err != nil {
		return nil, errors.Wrap(err, "inserting new user in DB")
	}
	return new_user, nil
}

func GetUser(user *models.User) error {
	db := database.GetDB()
	err := db.Model(user).
		Where("id = ?", user.ID).
		Column("id").
		First()
	return errors.Wrap(err, "check if user exists by ID")
}

func GetFullUser(user *models.User) error {
	db := database.GetDB()
	err := db.Model(user).
		Where("id = ?", user.ID).
		Column("user.*", "TendressesSent").
		Column("user.*", "TendressesReceived").
		Column("user.*", "Roles").
		Column("user.*", "Friends").
		Column("user.*", "Achievements").
		First()
	return errors.Wrap(err, "get full user")
}

func GetProfile(user *models.User) error {
	db := database.GetDB()
	err := db.Model(user).
		Where("id = ?", user.ID).
		Column("user.*", "Achievements").
		Column("user.*", "Friends").
		First()
	return errors.Wrap(err, "get profile of user")
}

func GetProfileByUsername(user *models.User, username string) error {
	db := database.GetDB()
	err := db.Model(user).
		Column("user.*", "Achievements").
		Column("user.*", "Friends").
		Where("username = ?", username).
		First()
	return errors.Wrap(err, "get profile of user")
}

func GetUserWithFriends(user *models.User) error {
	db := database.GetDB()
	err := db.Model(user).
		Where("id = ?", user.ID).
		Column("user.*", "Friends").
		First()
	return errors.Wrap(err, "get user and its friends")
}

func GetUserByUsername(user *models.User, username string) error {
	db := database.GetDB()
	err := db.Model(user).
		Where("username = ?", username).
		First()
	return errors.Wrap(err, "get user by username")
}

func AddRoleToUser(role *models.Role, user *models.User) error {
	db := database.GetDB()
	ur := new(models.UsersRoles)
	ur.RoleID = role.ID
	ur.UserID = user.ID
	err := db.Insert(ur)
	return errors.Wrap(err, "adding role to user")
}

func RemoveRoleFromUser(role *models.Role, user *models.User) error {
	db := database.GetDB()
	ur := new(models.UsersRoles)
	ur.RoleID = role.ID
	ur.UserID = user.ID
	err := db.Delete(ur)
	return errors.Wrap(err, "removing role from user")
}

func AddFriend(current_user *models.User, friend *models.User) error {
	db := database.GetDB()
	if err := GetUserWithFriends(current_user); err != nil {
		return err
	}
	for _, cu_friend := range current_user.Friends {
		if cu_friend.ID == friend.ID {
			return nil
		}
	}
	friendship := new(models.UsersFriends)
	friendship.UserID = current_user.ID
	friendship.FriendID = friend.ID
	err := db.Insert(friendship)
	return errors.Wrap(err, "creating new friendship with current_user")
}

func RemoveFriend(current_user *models.User, friend *models.User) error {
	db := database.GetDB()
	if err := GetUserWithFriends(current_user); err != nil {
		return err
	}
	for _, cu_friend := range current_user.Friends {
		if cu_friend.ID == friend.ID {
			friendship := new(models.UsersFriends)
			friendship.UserID = current_user.ID
			friendship.FriendID = friend.ID
			err := db.Delete(friendship)
			return errors.Wrap(err, "removing friendship")
		}
	}
	return nil
}

func UpdateUserAchievements(user *models.User, achievements []*models.Achievement, type_of string) error {
	if type_of != "send" && type_of != "receive" {
		return errors.New("wrong type, sender or receiver")
	}
	db := database.GetDB()
	for _, achievement := range achievements {
		if achievement.Type == type_of {
			ua := new(models.UsersAchievements)
			ua.UserID = user.ID
			ua.AchievementID = achievement.ID
			if err := GetOrCreateUserAchievement(ua); err != nil {
				log.Error(err)
				continue
			}
			if ua.Unlocked {
				continue
			}
			if ua.Score < achievement.Condition {
				ua.Score++
			} else {
				ua.Unlocked = true
			}
			if err := db.Update(ua); err != nil {
				log.Error(errors.Wrap(err, fmt.Sprintf("updating %q achievement", type_of)))
			}
		}
	}
	return nil
}

func GetOrCreateUserAchievement(ua *models.UsersAchievements) error {
	db := database.GetDB()
	if err := db.Model(ua).
		Where("achievement_id = ?", ua.AchievementID).
		Where("user_id = ?", ua.UserID).
		First(); err == nil {
		return nil
	}
	err := db.Insert(ua)
	return errors.Wrap(err, "creating UserAchievement after trying to get it")
}
