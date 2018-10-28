package models

import "fmt"

type Gif struct {
	tableName struct{} `sql:"gifs,alias:gif"`

	BaseModel
	Url        string       `json:"url",sql:"notnull,unique"`
	LameScore  int          `json:"lame_score,omitempty"`
	BlogID     int          `json:"blog_id,omitempty"`
	Blog       Blog         `json:"blog,omitempty"`
	Tags       []Tag        `pg:",many2many:gifs_tags,joinFK:tag_id" json:"tags,omitempty"`
	Tendresses []*Tendresse `json:"tendresses,omitempty"`
}

func (a Gif) String() string {
	return fmt.Sprintf("Gif<ID=%d Url=%q>", a.ID, a.Url)
}
