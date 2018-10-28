package models

import "fmt"

type Achievement struct {
	tableName struct{} `sql:"achievements,alias:achievement"`

	BaseModel
	Name      string `json:"name",sql:"unique,notnull"`
	Condition int    `json:"condition,omitempty",sql:",notnull"`
	Icon      string `json:"icon,omitempty"`
	Type      string `json:"type,omitempty",sql:",notnull"`
	Xp        int    `json:"xp,omitempty"`
	TagID     int    `json:"tag_id,omitempty"`
	Tag       *Tag   `json:"tag,omitempty"`
}

func (a Achievement) String() string {
	return fmt.Sprintf("Achievement<ID=%d, Name=%q>", a.ID, a.Name)
}
