# Assignment 2 — Cloud-Based Dashboard Service

> **Note for reviewers:** The service is fully deployed and live on our OpenStack VM.
> You do **not** need to clone or redeploy the project to test it.
> The GitLab repository is provided to verify code structure, testing, implementation, and deployment setup.

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

## Deployment

This project is containerized and runs on an OpenStack VM.

### Steps to deploy on OpenStack VM:

1. **Clone the repository**:
   ```bash
   git clone git@github.com:Amundfpl/cloud-dashboard-service.git
   cd Assignment-2
   ```

2. **(Optional)** If redeploying on a different VM, place the Firebase key:
    - Add your `firebase-key.json` inside the `credentials/` directory.
    -  This file is intentionally excluded from Git via `.gitignore`.

   > **Note:** The Firebase key is already mounted on our deployed VM. **Reviewers do not need access to it**. Only required if redeploying.

3. **Run the deployment script**:
   ```bash
   ./deploy.sh
   ```
> **Note:** If you get a **permission denied** error when trying to run the deploy script, make sure it’s executable by running:  
> `chmod +x deploy.sh`
> This makes the script executable so you can run it with `./deploy.sh`.


   This will:
    - Build the Docker image
    - Stop and remove any previous container
    - Mount credentials (if provided)
    - Expose the service on port `8080`

---

## Deployed Service URL

The service is live at:

```
add render link... or other deployment link
```

> Ensure port 8080 is open on your OpenStack VM.

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
│       ├── devops.yml
│       └── sync.yaml
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
├── deploy.sh                         # Deployment script
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

| Name    | Main Contributions                      |
|---------|------------------------------------------|
| Amund   | Caching logic, Dockerization, Deployment |
| Halvard | Webhook system, Testing, Service logic   |

We worked closely on most of the project. These areas highlight primary focus areas.

External data from:
- [REST Countries](https://restcountries.com/)
- [Open-Meteo](https://open-meteo.com/)
- [Frankfurter](https://www.frankfurter.app/)

---

## Notes

- Please ensure you do **not** commit your Firebase key file.
- All endpoints were tested using Go's built-in `httptest` package.
- External APIs were stubbed in tests to ensure no real requests were made.
- We implemented **advanced caching with purge**, full webhook triggering, PATCH/HEAD/DELETE support, and proper Docker deployment.
- The service follows RESTful principles and is easy to deploy via `deploy.sh` on any Docker-ready VM.

