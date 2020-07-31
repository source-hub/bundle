
package models

import (
	"time"
)

type Item_Catalog struct{
	Id uint `gorm:"primary_key AUTO_INCREMENT" json:"id" sql:"type:bigserial"`
	Item_id uint 	`json:"item_id" gorm:"type:bigserial REFERENCES items(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	Catalog_id uint `json:"catalog_id" gorm:"type:bigserial REFERENCES catalogs(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" sql:"index"`
}
func (Tc Item_Catalog) TableName() string {
	return "item_catalogs"
}
