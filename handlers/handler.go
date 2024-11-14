package handlers

import (
	"fmt"
	"net/http"

	"github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/db"
	"github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

var MaxAmount = 10000.00

func GetAccount(c *gin.Context) {
	var account models.Account
	var input struct {
		UserID uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err,
		}).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("user_id=?", input.UserID).First(&account).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("Cannot find user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Response := models.AccountResponse{
		AccountNumber: account.AccountNumber,
		Balance:       account.Balance,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account retrieved successfully",
		"Account": Response,
	})
}

func CreateAccount(c *gin.Context) {
	var account models.Account

	var input struct {
		AccountNumber string `json:"account_number"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Where("account_number=?", input.AccountNumber).First(&account).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": fmt.Sprintf("The account number %s aldready exists", input.AccountNumber)})
		return
	}

	if len(input.AccountNumber) != 12 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Length should be minimum 5"})
		return
	}

	account.AccountNumber = input.AccountNumber
	account.Balance = 0

	if err := db.DB.Create(&account).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error creating account")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.AccountResponse{
		AccountNumber: account.AccountNumber,
		Balance:       account.Balance,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account created successfully",
		"Account": response,
	})
}

func DepositAmount(c *gin.Context) {
	var account models.Account
	var input struct {
		UserID uint    `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.DB.Begin()

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, input.UserID).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("cannot find account")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if input.Amount < 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid amount"})
		return
	}

	if input.Amount > MaxAmount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Cannot deposit Above %.2f", MaxAmount)})
		return
	}

	account.Balance += input.Amount
	if err := tx.Save(&account).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("cannot save account")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	transaction := models.Transactions{
		AccountID:       input.UserID,
		Amount:          input.Amount,
		TransactionType: "Credit",
		Description:     "Credited amount to account",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("error creating transaction")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating transaction"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("Transaction commit failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}
	log.Info("Deposit Successfull")
	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit Successfull",
		"new_balance": account.Balance,
	})

}

func Withdraw(c *gin.Context) {
	var account models.Account
	var input struct {
		UserID uint    `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("Invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tx := db.DB.Begin()

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&account, input.UserID).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("Cannot find account")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if account.Balance < input.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient balance"})
		return
	}

	if input.Amount > MaxAmount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot withdraw over 10000"})
		return
	}

	if input.Amount < 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid amount"})
		return
	}

	account.Balance -= input.Amount
	if err := tx.Save(&account).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("Cannot udpate account")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot update account"})
		return
	}

	transaction := models.Transactions{
		AccountID:       input.UserID,
		Amount:          input.Amount,
		TransactionType: "Debit",
		Description:     "Debited amount from account",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("error creating transaction")
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating transaction"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"UserID": input.UserID,
			"error":  err.Error(),
		}).Error("error commiting transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error commiting transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Amount debited successfully",
		"new_balance": account.Balance,
	})
}

func TransferAmount(c *gin.Context) {
	var input struct {
		SenderAccountNumber   string  `json:"sender_account_number"`
		ReceiverAccountNumber string  `json:"receiver_account_number"`
		Amount                float64 `json:"amount"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transfer amount must be positive"})
		return
	}

	tx := db.DB.Begin()

	var sender, receiver models.Account

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("account_number = ?", input.SenderAccountNumber).First(&sender).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender account not found"})
		return
	}

	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("account_number = ?", input.ReceiverAccountNumber).First(&receiver).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Receiver account not found"})
		return
	}

	if sender.Balance < input.Amount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient funds"})
		return
	}

	sender.Balance -= input.Amount
	receiver.Balance += input.Amount

	if err := tx.Save(&sender).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sender's balance"})
		return
	}

	if err := tx.Save(&receiver).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update receiver's balance"})
		return
	}

	senderTransaction := models.Transactions{
		AccountID:       sender.ID,
		Amount:          -input.Amount,
		TransactionType: "Transfer Out",
		Description:     fmt.Sprintf("Transfer to %s", input.ReceiverAccountNumber),
	}
	if err := tx.Create(&senderTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sender transaction record"})
		return
	}

	receiverTransaction := models.Transactions{
		AccountID:       receiver.ID,
		Amount:          input.Amount,
		TransactionType: "Transfer In",
		Description:     fmt.Sprintf("Transfer from %s", input.SenderAccountNumber),
	}
	if err := tx.Create(&receiverTransaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create receiver transaction record"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction commit failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "Transfer successful",
		"sender_new_balance":   sender.Balance,
		"receiver_new_balance": receiver.Balance,
	})
}
