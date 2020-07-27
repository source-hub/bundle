package models

import (
	"time"
)

type Catalog struct {
	ID         uint       `gorm:"primary_key AUTO_INCREMENT" json:"id" sql:"type:bigserial"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" sql:"index"`
	Name string `json:"name"`
	Root string `json:"root"`
	Parent string `json:"parent"`
	Type string `json:"type"`
	Votes int `json:"votes"`
	Tags []string `json:"tags"`
	User_id uint `json:"user_id" sql:"type bigserial REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE"`
}

func (c *Catalog) TableName() string {
	return "catalogs"
}
