package models

import (
	"time"
)
type Tag struct{
	Name string `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" sql:"index"`
	Tag_Catalogs []Tag_Catalog `gorm:"foreignkey:tag_id"`
}
func (t *Tag)TableName() string{
	return "tags"
}