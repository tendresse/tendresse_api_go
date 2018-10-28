package models

import (
	"github.com/go-pg/pg/orm"
	"time"
)

type BaseModel struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,omitempty",sql:"default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty",sql:"default:now()"`
}

func (b *BaseModel) BeforeInsert(db orm.DB) error {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now()
	}
	return nil
}

func (b *BaseModel) BeforeUpdate(db orm.DB) error {
	if b.UpdatedAt.IsZero() {
		b.UpdatedAt = time.Now()
	}
	return nil
}
