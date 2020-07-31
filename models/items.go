package models

import (
	"time"
)

type Item struct {
	ID         uint       `gorm:"primary_key AUTO_INCREMENT" json:"id" sql:"type:bigserial"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" sql:"index"`
	Link string `json:"link"`
	Img string `json:"img"`
	Item_catalogs []Item_Catalog `gorm:"foreignkey:item_id"`
}

func (c *Item) TableName() string {
	return "items"
}
