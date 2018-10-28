package models

import "fmt"

type Role struct {
	tableName struct{} `sql:"roles,alias:role"`

	BaseModel
	Name        string `json:"name",sql:"unique,notnull"`
	Description string `json:"description,omitempty"`
}

func (a Role) String() string {
	return fmt.Sprintf("Role<ID=%d, Name=%q>", a.ID, a.Name)
}
