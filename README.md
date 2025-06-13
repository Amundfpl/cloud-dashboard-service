# Assignment 2 — Cloud-Based Dashboard Service

This project is a backend web service built in Go, originally developed as part of the Cloud Technologies course (PROG2005) at NTNU. It exposes a RESTful dashboard API with support for configurable dashboards, live data enrichment, caching, webhook notifications, and advanced HTTP methods (PATCH, HEAD, DELETE). The service is containerized using Docker and was initially deployed on an OpenStack VM.

> **This repository is a public-friendly refactored version of the original coursework project.**  
> The internal NTNU-hosted APIs (such as the course-provided country and currency services) have been replaced with fully external, publicly available APIs:
>
> - `REST Countries API` (https://restcountries.com)
> - `Open-Meteo API` (https://open-meteo.com)
> - `Frankfurter API` (https://www.frankfurter.app)
>
> Some internal logic and structure have been updated accordingly for clarity, reusability, and public deployment.

The original project was created by **Amund** and **Halvard** as part of the course submission. This repository represents a **fun spin-off and personal refactor** of that work, aimed at sharing the architecture and concepts openly on GitHub.

---

## External APIs Used

-  **REST Countries**  
  `https://restcountries.com/v3.1/alpha/{isoCode}`  
  Provides country metadata (capital, lat/lon, population, area, etc.)

- **Open-Meteo Weather API**  
  `https://api.open-meteo.com/v1/forecast?...`  
  Provides current temperature and precipitation data

- **Frankfurter Currency API**  
  `https://api.frankfurter.app/latest?from=EUR&to=USD,NOK`  
  Provides exchange rates between currency pairs

---

## Deployed Service URL

The service is live at:

```
https://cloud-dashboard-service.onrender.com/
```
> **Note for reviewers:** Render is a free hosting service that puts the app to sleep after 15 minutes of inactivity. 
> So the first request may take a few seconds to wake up. Subsequent requests will be faster.

---

## API Endpoints

The service supports the following RESTful endpoints:

### `/dashboard/v1/registrations/`

Register, view, update, delete, and manage dashboard configurations.

#### POST - Register a new dashboard
```http
POST /dashboard/v1/registrations/
Content-Type: application/json
```
**Request body example:**
```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": true,
    "targetCurrencies": ["EUR", "USD"]
  }
}
```
**Response:**
```json
{
  "id": "abc123",
  "lastChange": "20250407 15:30"
}
```

#### GET - View specific configuration
```http
GET /dashboard/v1/registrations/{id}
```
**Response:**
```json
{
  "id": "abc123",
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": true,
    "precipitation": true,
    "capital": true,
    "coordinates": true,
    "population": true,
    "area": true,
    "targetCurrencies": ["EUR", "USD"]
  },
  "lastChange": "20250407 15:30"
}
```

#### PATCH / PUT / DELETE / HEAD supported for advanced config updates

---

### `/dashboard/v1/dashboards/{id}`

Get a fully enriched dashboard based on the config.
```http
GET /dashboard/v1/dashboards/{id}
```
**Response:**
```json
{
  "country": "Norway",
  "isoCode": "NO",
  "features": {
    "temperature": 4.5,
    "precipitation": 0.7,
    "capital": "Oslo",
    "coordinates": {"latitude": 59.9, "longitude": 10.8},
    "population": 5400000,
    "area": 385207,
    "targetCurrencies": {
      "EUR": 0.087,
      "USD": 0.092
    }
  },
  "lastRetrieval": "20250407 16:00"
}
```

---

### `/dashboard/v1/notifications/`

#### Supported Webhook Events:
- `REGISTER` — When a new dashboard is created
- `CHANGE` — When a dashboard is updated via PUT
- `PATCH` — When partially updated
- `DELETE` — When a dashboard is deleted
- `INVOKE` — When a dashboard is accessed (GET)
- `LOW_TEMP` — **When the temperature is below 0°C during dashboard enrichment**


#### POST - Register a webhook
```http
POST /dashboard/v1/notifications/
Content-Type: application/json
```
**Body example:**
```json
{
  "url": "https://webhook.site/your-hook",
  "country": "NO",
  "event": "INVOKE"
}
```
**Response:**
```json
{
  "id": "webhook123"
}
```

#### GET/DELETE by ID supported.

Webhook POST payloads:
```json
{
  "id": "webhook123",
  "country": "NO",
  "event": "INVOKE",
  "time": "20250407 16:30"
}
```

---

### `/dashboard/v1/status/`

Returns current system health, API dependencies, and uptime.
```json
{
  "countries_api": 200,
  "meteo_api": 200,
  "currency_api": 200,
  "notification_db": 200,
  "webhooks": 4,
  "version": "v1",
  "uptime": 3021
}
```

---

## Example `curl` Commands

Register dashboard:
```bash
curl -X POST http://localhost:8080/dashboard/v1/registrations/ \
  -H "Content-Type: application/json" \
  -d '{"country":"Norway","isoCode":"NO","features":{"capital":true}}'
```

Register webhook:
```bash
curl -X POST http://localhost:8080/dashboard/v1/notifications/ \
  -H "Content-Type: application/json" \
  -d '{"url":"https://webhook.site/abc","event":"INVOKE"}'
```

---

## Project Structure

```plaintext
Assignment-2/
├── .github/
│   └── workflows/
│       └── devops.yaml
├── cache/
│   ├── cache_autoPurge.go
│   ├── cache_autoPurge_test.go
│   ├── cache_keys.go
│   ├── cache_keys_test.go
│   ├── cache_purge.go
│   ├── cache_purge_test.go
│   ├── cache_store.go
│   └── cache_store_test.go
├── cmd/
│   └── main.go
├── credentials/
│   ├── firebase-key.json              # Not committed to repo
│   └── test-serviceAccountKey.json    # For local testing
├── db/
│   ├── firebase.go
│   ├── repository.go
│   └── webhook_db.go
├── handlers/
│   ├── dashboard_handler.go
│   ├── dashboard_handler_test.go
│   ├── notification_handler.go
│   ├── notification_handler_test.go
│   ├── registration_handler.go
│   ├── registration_handler_test.go
│   ├── service_handler.go
│   └── service_handler_test.go
├── httpclient/
│   └── httpClient.go
├── server/
│   ├── dbInit.go
│   ├── router.go
│   └── server.go
├── services/
│   ├── dashboard_service.go
│   ├── dashboard_service_test.go
│   ├── enrichment_service.go
│   ├── enrichment_service_test.go
│   ├── notification_service.go
│   ├── notification_service_test.go
│   ├── registration_service.go
│   ├── registration_service_test.go
│   ├── status_service.go
│   └── status_service_test.go
├── static/
│   └── index.html                     # Homepage file served from "/"
├── testsetup/
│   └── setup.go                       # Helpers for setting up mocks, test env
├── utils/
│   ├── config.go
│   ├── dashboardMessages.go
│   ├── models.go
│   ├── responseUtil.go
│   └── util.go
├── .gitignore
├── Dockerfile                        # Multi-stage Docker build
├── go.mod
├── go.sum
└── README.md

```

---

## Running Tests

This project uses Go’s built-in testing framework. Below are the most common commands to run and debug tests.

```bash
# Run all tests in the project
go test ./...

# Run tests in a specific package (e.g., services)
go test ./services

# Run tests with verbose output (show details of each test)
go test -v ./...

# Run a specific test function by name (e.g., TestGetEnrichedDashboards)
go test -v ./services -run TestGetEnrichedDashboards

# Clean test cache (helpful if test results seem outdated)
go clean -testcache

# Tidy up dependencies (optional, for cleanup)
go mod tidy
```
---

## Tech Used

- Go 1.22+
- Firebase Firestore
- Docker
- REST APIs (Countries, Currency, Weather)
- OpenStack VM
- GitLab CI/CD

---

## Credits

External data from:
- [REST Countries](https://restcountries.com/)
- [Open-Meteo](https://open-meteo.com/)
- [Frankfurter](https://www.frankfurter.app/)

---

## Notes

- All endpoints were tested using Go's built-in `httptest` package.
- External APIs were stubbed in tests to ensure no real requests were made.
- We implemented **advanced caching with purge**, full webhook triggering, PATCH/HEAD/DELETE support, and proper Docker deployment.
- The service follows RESTful principles.

