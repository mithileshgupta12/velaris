package models

import "time"

type User struct {
	Id        int64
	Name      string    `xorm:"NOT NULL"`
	Email     string    `xorm:"NOT NULL UNIQUE"`
	Password  string    `xorm:"NOT NULL" json:"-"`
	Boards    []*Board  `xorm:"-"`
	CreatedAt time.Time `xorm:"timestamptz NOT NULL created"`
	UpdatedAt time.Time `xorm:"timestamptz NOT NULL updated"`
}
