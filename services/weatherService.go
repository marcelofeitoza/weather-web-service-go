package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/marcelofeitoza/weather-web-service-go/models"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/url"
	"time"
)

func fetchLatLong(city string, db *sqlx.DB, rdb *redis.Client) (*models.LatLong, error) {
	var latlong models.LatLong
	endpoint := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", url.QueryEscape(city))
	ctx := context.Background()

	latlongJSON, err := rdb.Get(ctx, fmt.Sprintf("lat-long-%s", city)).Result()
	if latlongJSON != "" {
		err = json.Unmarshal([]byte(latlongJSON), &latlong)
		if latlong.Latitude != 0 && latlong.Longitude != 0 {
			return &latlong, nil
		}
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making request to geocoding API: %w", err)
	}
	defer resp.Body.Close()

	var response models.GeoResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("error decoding response from geocoding API: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no results found for city: %s", city)
	}

	latlong = response.Results[0]

	_, err = db.Exec("INSERT INTO cities (name, lat, long) VALUES ($1, $2, $3)", city, latlong.Latitude, latlong.Longitude)

	latlongJSONBytes, err := json.Marshal(latlong)
	if err != nil {
		return nil, fmt.Errorf("error marshalling latlong: %w", err)
	}

	latlongJSON = string(latlongJSONBytes)
	if err != nil {
		return nil, fmt.Errorf("error marshalling latlong: %w", err)
	}

	err = rdb.Set(ctx, fmt.Sprintf("lat-long-%s", city), string(latlongJSON), 60*time.Minute).Err()

	return &latlong, nil
}

func GetLatLong(city string, db *sqlx.DB, rdb *redis.Client) (*models.LatLong, error) {
	var latlong models.LatLong
	ctx := context.Background()

	latlongStr, err := rdb.Get(ctx, fmt.Sprintf("lat-long-%s", city)).Result()
	if latlongStr != "" {
		_ = json.Unmarshal([]byte(latlongStr), &latlong)
		return &latlong, nil
	}

	err = db.Get(
		&latlong,
		"SELECT lat AS latitude, long AS longitude FROM cities WHERE name = $1",
		city,
	)
	if err != nil {
		latlongPtr, err := fetchLatLong(city, db, rdb)
		if err != nil {
			return nil, fmt.Errorf("error fetching latlong: %w", err)
		}
		latlong = *latlongPtr
		return &latlong, nil
	}

	return &latlong, nil
}

func GetWeather(latLong models.LatLong, rdb *redis.Client) (*models.WeatherResponse, error) {
	var weatherResponse models.WeatherResponse
	endpoint := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&hourly=temperature_2m", latLong.Latitude, latLong.Longitude)

	ctx := context.Background()
	weatherResponseJSON, err := rdb.Get(ctx, fmt.Sprintf("weather-%f-%f", latLong.Latitude, latLong.Longitude)).Result()
	if err == nil {
		err = json.Unmarshal([]byte(weatherResponseJSON), &weatherResponse)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling data from Redis: %w", err)
		}
		return &weatherResponse, nil
	}

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making request to Weather API: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	bytes, _ := json.Marshal(weatherResponse)
	weatherResponseJSON = string(bytes)
	err = rdb.Set(ctx, fmt.Sprintf("weather-%f-%f", latLong.Latitude, latLong.Longitude), weatherResponseJSON, 60*time.Second).Err()

	return &weatherResponse, nil
}
