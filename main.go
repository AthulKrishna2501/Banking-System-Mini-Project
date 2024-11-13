package main

import (
	db "github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/db"
	"github.com/AthulKrishna2501/Banking-System-Mini-Project-.git/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	db.InitDatabase()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	router := gin.Default()

	router.GET("/get-account", routes.GetAccount)

}
