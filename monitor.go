package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func notifyWebhook(url string) {
	res, err := httpClient.Get(url)

	if err == nil {
		res.Body.Close()
	}
}

func monitor(id uuid.UUID, site string, webhook string, interval int) {
	for {
		// Check site existence
		var siteToMonitor ToMonitor
		if err := DB.First(&siteToMonitor, "id = ?", id).Error; err != nil {
			return
		}

		// Check site
		notifyWebhook(webhook + "/start")

		res, err := httpClient.Get(site)

		if err != nil {
			fmt.Println(site + " is down")
			notifyWebhook(webhook + "/fail")
		} else if res.StatusCode >= 500 {
			fmt.Println(site + " is down")
			notifyWebhook(webhook + "/fail")
			res.Body.Close()
		} else {
			fmt.Println(site + " is up")
			notifyWebhook(webhook)
			res.Body.Close()
		}

		// Wait for interval (min)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func fireMonitors() {
	var rows []ToMonitor
	res := DB.Find(&rows)

	if res.Error != nil {
		log.Fatal("Failed to fetch existing sites")
	}

	for _, row := range rows {
		go monitor(row.ID, row.Site, row.Webhook, row.Interval)
	}
}
