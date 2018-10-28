package models

import "fmt"

type Tag struct {
	tableName struct{} `sql:"tags,alias:tag"`

	BaseModel
	Name         string         `json:"name",sql:"unique,notnull"`
	Achievements []*Achievement `json:"achievements,omitempty"`
	Gifs         []Gif          `pg:",many2many:gifs_tags,joinFK:gif_id" json:"gifs,omitempty"`
}

func (a Tag) String() string {
	return fmt.Sprintf("Tag<ID=%d Name=%q>", a.ID, a.Name)
}

// func (tag *Tag) get_random_tag() {
// 	db.Order("RANDOM()").Find(&tag)
// }
