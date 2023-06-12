package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"log"

	"github.com/espher/GoLang-API-REST/db"
	"github.com/espher/GoLang-API-REST/models"
	"github.com/espher/GoLang-API-REST/routes"
	"github.com/espher/GoLang-API-REST/sqs"
)

func main() {
	db.ConnectDB()

	err := db.DB.AutoMigrate(&models.User{}, &models.Post{})
	if err != nil {
		log.Fatal("Failed to run databases:", err)
	}

	sqs.GetQueueList()

	router := gin.Default()
	routes.SetupRoutes(router)
	http.ListenAndServe(":9898", router)
	//router.Run("localhost:9898")
}
