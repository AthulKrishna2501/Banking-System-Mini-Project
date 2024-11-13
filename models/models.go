package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	AccountNumber string `gorm:"uniqueIndex"`
	UserID        uint
	Balance       float64 `gorm:"not null"`
}

type Transations struct {
	gorm.Model
	AccountID       uint
	Amount          float64
	TransactionType string
	Description     string
}
