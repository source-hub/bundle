package models

import (
	"time"
)

type User struct {
	ID         uint       `gorm:"primary_key AUTO_INCREMENT" json:"id" sql:"type:bigserial"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"-" sql:"index"`
	Password   string     `json:"password"`
	Zip        uint       `json:"zip" gorm:"not null;"`
	First_name string     `json:"first_name"`
	Last_name  string     `json:"last_name"`
	Guid  string
	Email string `json:"email" gorm:"not null;unique"`
	Catalogs []Catalog	`gorm:"foreignkey:user_id"` 
	//
}

func (U *User) TableName() string {
	return "users"
}
