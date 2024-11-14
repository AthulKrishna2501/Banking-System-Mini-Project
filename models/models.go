package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	AccountNumber string
	UserID        uint    `gorm:"autoIncrement"`
	Balance       float64 `gorm:"not null"`
}

type Transactions struct {
	gorm.Model
	AccountID       uint
	Amount          float64
	TransactionType string
	Description     string
}

type AccountResponse struct {
	AccountNumber string  `json:"account_number"`
	Balance        float64 `json:"balance"`
}
