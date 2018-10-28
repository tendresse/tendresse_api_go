package models

import "fmt"

type Blog struct {
	tableName struct{} `sql:"blogs,alias:blog"`

	BaseModel
	Url         string `json:"url",sql:"unique,notnull"`
	Description string `json:"description"`
	Gifs        []*Gif `json:"gifs,omitempty"`
}

func (a Blog) String() string {
	return fmt.Sprintf("Blog<ID=%d Description=%q>", a.ID, a.Description)
}
