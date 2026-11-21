package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getSites(ctx *gin.Context) {
	var rows []ToMonitor
	res := DB.Find(&rows)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on fetching"})
		return
	}

	ctx.JSON(http.StatusOK, rows)
}

func createNewSite(ctx *gin.Context) {
	var siteToMonitor ToMonitor
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

	res := DB.Create(&siteToMonitor)

	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on creation"})
		return
	}

	go monitor(siteToMonitor.ID, siteToMonitor.Site, siteToMonitor.Webhook, siteToMonitor.Interval)

	ctx.JSON(http.StatusCreated, siteToMonitor)
}

func deleteSite(ctx *gin.Context) {
	var siteToMonitor ToMonitor
	if err := DB.First(&siteToMonitor, "id = ?", ctx.Param("id")).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	res := DB.Delete(&siteToMonitor)
	if res.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on deletion"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
