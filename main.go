package main

import (
	db "github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/db"
	"github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/handlers"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	db.InitDatabase()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	router := gin.Default()

	router.GET("/get-account", handlers.GetAccount)
	router.POST("/create-account", handlers.CreateAccount)
	router.POST("/add-amount", handlers.DepositAmount)
	router.PUT("/withdraw-amount", handlers.Withdraw)
	router.POST("/transfer-amount",handlers.TransferAmount)

	router.Run(":3000")
}
