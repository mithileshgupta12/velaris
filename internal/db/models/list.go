package models

import "time"

type List struct {
	Id        int64     `json:"id"`
	Name      string    `xorm:"NOT NULL" json:"name"`
	BoardId   int64     `xorm:"INDEX NOT NULL" json:"board_id"`
	Board     *Board    `xorm:"-" json:"board"`
	Position  int       `xorm:"NOT NULL" json:"position"`
	CreatedAt time.Time `xorm:"NOT NULL created" json:"created_at"`
	UpdatedAt time.Time `xorm:"NOT NULL updated" json:"updated_at"`
}

func (l *List) TableName() string {
	return "lists"
}
