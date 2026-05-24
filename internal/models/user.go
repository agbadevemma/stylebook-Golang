package models

type User struct {
	Base
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);unique;not null"`
	Phone    string `json:"phone" gorm:"type:varchar(255);unique;not null"`
	Password string `json:"password" gorm:"type:text;not null"`
	Role     string `json:"role" gorm:"type:enum('barber', 'admin');default:'barber';not null"`
}
