package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const apiUrl = "https://api.open-meteo.com/v1/forecast?latitude=56.8575&longitude=60.6125&current=temperature_2m"

type ResponseWeather struct {
	Latitude             float64      `json:"latitude"`
	Longitude            float64      `json:"longitude"`
	GenerationtimeMS     float64      `json:"generationtime_ms"`
	UTCOffsetSeconds     int64        `json:"utc_offset_seconds"`
	Timezone             string       `json:"timezone"`
	TimezoneAbbreviation string       `json:"timezone_abbreviation"`
	Elevation            float64      `json:"elevation"`
	CurrentUnits         CurrentUnits `json:"current_units"`
	Current              Current      `json:"current"`
}

type Current struct {
	Time          string  `json:"time"`
	Interval      int64   `json:"interval"`
	Temperature2M float64 `json:"temperature_2m"`
}

type CurrentUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature2M string `json:"temperature_2m"`
}

var resp ResponseWeather

func init() {
	getTemp()
}

func main() {
	go func() {
		for range time.Tick(60 * time.Second) {
			getTemp()
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		preparedData, jsonErr := json.Marshal(resp.Current.Temperature2M)

		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(preparedData)
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}

func getTemp() {
	data, err := http.Get(apiUrl)

	if err != nil {
		getTemp()
	}

	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &resp)
}
