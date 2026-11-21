package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func main() {
	connectToDatabase()
	if err := DB.AutoMigrate(&ToMonitor{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// fire monitoring existing sites
	go fireMonitors()

	// Define Routes
	router := gin.Default()

	router.GET("/sites", getSites)
	router.POST("/sites/new", createNewSite)
	router.DELETE("/sites/:id", deleteSite)

	log.Println("Server running on localhost:8080")
	router.Run("localhost:8080")
}
