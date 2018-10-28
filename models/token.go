package models

import "fmt"

type Token struct {
	tableName struct{} `sql:"tokens,alias:token"`

	BaseModel
	Hash   string `json:"-",sql:"unique,notnull"`
	UserID int    `json:"user_id",sql:",notnull"`
	User   *User  `json:"user,omitempty"`
}

func (a Token) String() string {
	return fmt.Sprintf("Token<ID=%d>", a.ID)
}
