package monitor

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/yahrour/uptime-monitor-api/internal/database"
	"github.com/yahrour/uptime-monitor-api/internal/models"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func NotifyWebhook(url string) {
	res, err := httpClient.Get(url)

	if err == nil {
		res.Body.Close()
	}
}

func Monitor(id uuid.UUID, site string, webhook string, interval int) {
	for {
		// Check site existence
		var siteToMonitor models.ToMonitor
		if err := database.DB.First(&siteToMonitor, "id = ?", id).Error; err != nil {
			return
		}

		// Check site
		NotifyWebhook(webhook + "/start")

		res, err := httpClient.Get(site)

		if err != nil {
			fmt.Println(site + " is down")
			NotifyWebhook(webhook + "/fail")
		} else if res.StatusCode >= 500 {
			fmt.Println(site + " is down")
			NotifyWebhook(webhook + "/fail")
			res.Body.Close()
		} else {
			fmt.Println(site + " is up")
			NotifyWebhook(webhook)
			res.Body.Close()
		}

		// Wait for interval (min)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func FireMonitors() {
	var rows []models.ToMonitor
	res := database.DB.Find(&rows)

	if res.Error != nil {
		log.Fatal("Failed to fetch existing sites")
	}

	for _, row := range rows {
		go Monitor(row.ID, row.Site, row.Webhook, row.Interval)
	}
}
