package models

import (
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	tableName struct{} `sql:"users,alias:user"`

	BaseModel
	Username           string        `json:"username",sql:"unique,notnull"`
	Passhash           string        `json:"-"`
	Email              string        `json:"-",sql:"unique,notnull"`
	Premium            bool          `json:"premium,omitempty"`
	Tokens             []*Token      `json:"tokens,omitempty"`
	TendressesReceived []*Tendresse  `pg:",fk:Receiver" json:"tendresses_received,omitempty" `
	TendressesSent     []*Tendresse  `pg:",fk:Sender" json:"tendresses_sent,omitempty"`
	Achievements       []Achievement `pg:",many2many:users_achievements,joinFK:achievement_id" json:"achievements,omitempty"`
	Friends            []User        `pg:",many2many:users_friends,joinFK:friend_id" json:"friends,omitempty"`
	Roles              []Role        `pg:",many2many:users_roles,joinFK:role_id" json:"roles,omitempty"`
}

func (a User) String() string {
	return fmt.Sprintf("User<ID=%d Username=%q>", a.ID, a.Username)
}

func (user *User) SetPassword(password string) error {
	encryptedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return errors.Wrap(err, "encrypting the user password")
	}
	user.Passhash = string(encryptedBytes[:])
	return nil
}

func (user *User) VerifyPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Passhash), []byte(password))
	if err != nil {
		return errors.Wrap(err, "comparing password during login")
	}
	return nil
}
