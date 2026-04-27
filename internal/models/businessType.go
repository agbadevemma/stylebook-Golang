package models

type BusinessType struct {
	Base
	Name string `gorm:"unique;not null" json:"name"`
}
