package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string    `json:"name" validate:"required"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password" validate:"required,min=8,max=100"`
	Expenses []Expense `gorm:"foreignKey:UserID"`
	Salt     string
}

type Expense struct {
	gorm.Model
	UserID   uint
	Title    string  `json:"title" validate:"required"`
	Amount   float64 `json:"amount" validate:"required,min=0"`
	Category string  `json:"category"`
}
