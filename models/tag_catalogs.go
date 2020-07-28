package models

import (
	"time"
)

type Tag_Catalog struct{
	Id uint `gorm:"primary_key AUTO_INCREMENT" json:"id" sql:"type:bigserial"`
	Tag_id uint 	`json:"tag_id" gorm:"type:bigserial REFERENCES tags(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Catalog_id uint `json:"catalog_id" gorm:"type:bigserial REFERENCES catalogs(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" sql:"index"`
}
func (Tc Tag_Catalog) TableName() string {
	return "tag_catalogs"
}
