package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yahrour/uptime-monitor-api/internal/database"
	"github.com/yahrour/uptime-monitor-api/internal/models"
	"github.com/yahrour/uptime-monitor-api/internal/monitor"
)

func GetSites(ctx *gin.Context) {
	var rows []models.ToMonitor
	res := database.DB.Find(&rows)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on fetching"})
		return
	}

	ctx.JSON(http.StatusOK, rows)
}

func CreateNewSite(ctx *gin.Context) {
	var siteToMonitor models.ToMonitor
	if err := ctx.ShouldBindJSON(&siteToMonitor); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if siteToMonitor.Interval < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Interval must be 1 (min) or more"})
		return
	}

	// http:// https://
	if len(siteToMonitor.Site) < 7 || len(siteToMonitor.Webhook) < 7 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid site/webhook url(s)"})
		return
	}

	res := database.DB.Create(&siteToMonitor)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on creation"})
		return
	}

	go monitor.Monitor(siteToMonitor.ID, siteToMonitor.Site, siteToMonitor.Webhook, siteToMonitor.Interval)

	ctx.JSON(http.StatusCreated, siteToMonitor)
}

func DeleteSite(ctx *gin.Context) {
	var siteToMonitor models.ToMonitor
	if err := database.DB.First(&siteToMonitor, "id = ?", ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	res := database.DB.Delete(&siteToMonitor)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on deletion"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
