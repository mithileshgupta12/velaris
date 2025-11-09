package models

import "time"

type Board struct {
	Id          int64     `json:"id"`
	Name        string    `xorm:"NOT NULL" json:"name"`
	Description *string   `xorm:"TEXT" json:"description"`
	UserId      int64     `xorm:"INDEX" json:"user_id"`
	User        *User     `xorm:"-" json:"user"`
	CreatedAt   time.Time `xorm:"timestamptz NOT NULL created" json:"created_at"`
	UpdatedAt   time.Time `xorm:"timestamptz NOT NULL updated" json:"updated_at"`
}

func (b *Board) TableName() string {
	return "boards"
}
