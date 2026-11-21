package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/yahrour/uptime-monitor-api/internal/database"
	"github.com/yahrour/uptime-monitor-api/internal/handlers"
	"github.com/yahrour/uptime-monitor-api/internal/models"
	"github.com/yahrour/uptime-monitor-api/internal/monitor"
)

func main() {
	database.ConnectToDatabase()
	if err := database.DB.AutoMigrate(&models.ToMonitor{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	// fire monitoring existing sites
	go monitor.FireMonitors()

	// Define Routes
	router := gin.Default()

	router.GET("/sites", handlers.GetSites)
	router.POST("/sites/new", handlers.CreateNewSite)
	router.DELETE("/sites/:id", handlers.DeleteSite)

	log.Println("Server running on localhost:8080")
	router.Run("localhost:8080")
}
