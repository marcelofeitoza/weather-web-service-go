# Go Weather Web Service

This is an implementation of a weather web service using Go, inspired by the [Rust Weather Web Service](https://github.com/marcelofeitoza/weather-web-service). It uses the `sqlx` package for PostgreSQL database storage and `go-redis` for the Redis caching system.

The service stores the last cities queried for weather checks in a PostgreSQL database, resulting in consultation times of around 1.7 seconds. To improve request times, Redis is used, which reduces consultation times to between 4 and 10 milliseconds.

---

## Running

To run the application, you will need [Go](https://golang.org/dl/) and [Docker](https://docs.docker.com/engine/install/) installed:

- Start Postgres and Redis:

  ```sh
  docker compose up -d # -d is for detached mode- terminal will be free
  ```

After running, the application will connect to Postgres and Redis and will be ready to receive requests at [http://localhost:8080](http://localhost:8080).

---

## Development

- Comment out this part of the `docker-compose.yml` file:
  ```md
    #  Comment for local development
    weather-web-service-go:
      build: ./
      ports:
        - "8080:8080"
      depends_on:
        - postgres
        - redis
      environment:
        DATABASE_URL: "postgres://forecast:forecast@forecast_postgres_go:54320/forecast?sslmode=disable"
        REDIS_URL: "redis://forecast_redis_go:63790"
    # Comment for local development
  ```

- Run the database and redis:

  ```sh
  docker compose up -d
  ```

- Run the application:

  ```sh
  go run main.go
  ```

---

## Testing

- Run the tests:

  ```sh
  go test ./...
  ```

---

To make a request for a city, simply access [http://localhost:8080/weather/your-city-name](http://localhost:8080/weather/your-city-name) and replace `your-city-name` with the name of the city you want to fetch. This will result in a response like this:

```json
{
  "duration": "893.557667ms",
  "city": "São Paulo",
  "forecasts": [
    {
      "date": "2024-02-16T00:00",
      "temperature": "19.3"
    },
    ...
  ]
}
```

In this example, this was the chosen city: [http://localhost:8080/weather/S%C3%A3o%20Paulo](http://localhost:8080/weather/S%C3%A3o%20Paulo).