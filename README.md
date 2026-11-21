# uptime-monitor-api

A lightweight uptime monitoring REST API built with Go, Gin, and SQLite, designed to work with [Healthchecks.io](https://healthchecks.io). It checks your sites at set intervals and reports the results to Healthchecks.io ping URLs, so you get alerts when a site goes down.

## Features

- Add, list, and delete monitored sites through a JSON API
- Configurable check interval (in minutes) per site
- Reports to Healthchecks.io using its start, success, and fail ping endpoints
- Monitors resume automatically after a restart
- Single-file SQLite database, no external setup

## Tech Stack

Go, [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io), SQLite

## Getting Started

### Prerequisites

- Go (see the version in `go.mod`)
- A [Healthchecks.io](https://healthchecks.io) account

### Run

```bash
git clone https://github.com/yahrour/uptime-monitor-api.git
cd uptime-monitor-api
go run .
```

The server listens on `localhost:8080` and creates an `app.db` SQLite file in the working directory on first run.

## Healthchecks.io Setup

1. Create a check in [Healthchecks.io](https://healthchecks.io).
2. Set its **Period** to the `interval` you'll use (in minutes), and give it a grace time of at least a minute so small delays don't trigger false alerts.
3. Copy the check's ping URL (e.g. `https://hc-ping.com/your-uuid`).
4. Register the site, using that ping URL as `webhook`:

   ```bash
   curl -X POST http://localhost:8080/sites/new \
     -H "Content-Type: application/json" \
     -d '{"site": "https://example.com", "webhook": "https://hc-ping.com/your-uuid", "interval": 5}'
   ```

Healthchecks.io then alerts you (email, Slack, etc.) if the pings stop or a `/fail` signal arrives.

## API Reference

| Method   | Endpoint     | Description              | Success          |
| -------- | ------------ | ------------------------ | ---------------- |
| `GET`    | `/sites`     | List all monitored sites | `200 OK`         |
| `POST`   | `/sites/new` | Add a site               | `201 Created`    |
| `DELETE` | `/sites/:id` | Delete a site            | `204 No Content` |

### Add a site

| Field      | Type    | Description                                   |
| ---------- | ------- | --------------------------------------------- |
| `site`     | string  | URL of the site to monitor                    |
| `webhook`  | string  | Healthchecks.io ping URL                      |
| `interval` | integer | Check interval in minutes (must be 1 or more) |

```bash
curl -X POST http://localhost:8080/sites/new \
  -H "Content-Type: application/json" \
  -d '{"site": "https://example.com", "webhook": "https://hc-ping.com/your-uuid", "interval": 5}'
```

## How It Works

Each monitored site runs in its own goroutine. It checks the site right away, then repeats every `interval` minutes:

1. Sends a start signal to `<ping-url>/start`
2. Requests the site (5 second timeout)
3. Sends a success ping to `<ping-url>` if the site is up, or a fail signal to `<ping-url>/fail` if the request failed or returned a 5xx status

The start signal lets Healthchecks.io measure how long each check takes. A fail signal marks the check as down immediately, without waiting for the grace time.

On startup, all sites stored in the database are loaded and their monitors are started again. Deleting a site stops its monitor the next time it wakes up.

If this service itself stops, the pings stop too, and Healthchecks.io alerts you once the period and grace time pass, so the monitor is monitored as well.

## Project Structure

```
.
├── main.go                  # Entry point, routes
└── internal/
    ├── database/            # SQLite connection
    ├── models/              # Data models
    ├── handlers/            # HTTP handlers
    └── monitor/             # Background monitoring logic
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
