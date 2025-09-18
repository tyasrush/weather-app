# WEATHER 

## DESCRIPTION
This is a weather application that provides weather forecasts for various locations.

## HOW-TO
1. You need to set env first by run `source env.sh`, you can copy existing file in project.
1. Setting up database resource by run `make docker-compose-up` or run manual docker command `docker compose up -d`.
1. Finally, you can run the app, rest api -> `make run-api` and worker -> `make run-worker`.

## ASSUMPTIONS
1. User will see current weather based on location they pick via endpoint `GET /weathers`, and also list of forecast data.
1. Current time weather will decide on filtering data from database, pick the one that near with time now.
1. Data collection will happen in endpoint `POST /weathers/sync` and worker (by running `make run-worker`).

## TRADE OFFS
1. When external weather api down, it make our app can't update the data, and stuck. Solution: find other api as a backup, crawling data from other sources, etc.
1. Potential overheat when running worker and external  weather api got issues, need to handle it.
1. If something happen to worker and make it stop work, there is no retry to make worker up, since worker still very simple.

## IMPROVEMENTS
Due to limited time, here are some improvements note.
1. Endpoint `POST /weathers/sync` slow, need to improve it's performance by implement concurrent process with goroutine or other async method.
1. Batching process on worker to avoid overheat CPU when insert to database.
1. Implement fallback when external api server down or failed.

## ENDPOINTS

### Locations
- GET /api/v1/locations - Get all locations
- POST /api/v1/locations - Create a new location

### Weather
- POST /api/v1/weathers/sync - Sync weather data
- GET /api/v1/weathers - Get weather data for a location

## COMMANDS

### Build and Run
- `make build-api` - Build the API server
- `make build-worker` - Build the weather sync worker
- `make run-api` - Build and run the API server
- `make run-worker` - Build and run the weather sync worker

### Testing
- `make test` - Run all tests

### Development
- `make deps` - Install dependencies
- `make mock` - Generate mocks
- `make clean` - Clean build artifacts
- `make docker-compose-up` - Start docker services
- `make docker-compose-down` - Stop docker services

## CONFIGURATION
- `PORT` - Application port (default: 8080)
- `MYSQL_DSN` - MySQL Data Source Name (default: admin:admin@tcp(localhost:3306)/weather-db?charset=utf8mb4&parseTime=true&loc=Local)
- `REDIS_ADDR` - Redis ip address (default: locatlho:6379)
- `REDIS_PASSWORD` - Redis password
- `WEATHER_API_BASE_URL` - Weather API Base URL (default: https://api.weatherapi.com/v1)
- `WEATHER_API_KEY` - Weather API Key to fetch data, generate apikey from your account here https://www.weatherapi.com/
- `BACKOFF_MAX_RETRIES` - Max retries for exponential backoff (default: 3)
- `BACKOFF_BASE_DELAY` - Initial wait duration for exponential backoff (default: 200ms)
- `BACKOFF_MAX_DELAY` - Max wait duration for exponential backoff (default: 5sec)
- `WORKER_PERIOD` - Period between sync operations in time duration type (default: 15min)
- `WORKER_LIMIT` - Maximum number of locations to sync (default: 10)

