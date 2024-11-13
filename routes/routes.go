package routes

import (
	"net/http"

	"github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/db"
	"github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var MaxAmount = 10000

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

	c.JSON(http.StatusOK, gin.H{
		"message": "Account retrieved successfullt",
		"Account": account,
	})
}

func CreateAccount(c *gin.Context) {
	var account models.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error binding JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account.Balance = 0

	if err := db.DB.Create(&account).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("error creating account")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account created succesfully",
		"Account": account,
	})
}

