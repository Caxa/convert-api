//http://localhost:8000/convert?from=centimeters&to=meters&value=150

package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func centimetersToMeters(centimeters float64) float64 {
	return centimeters / 100
}

func metersToCentimeters(meters float64) float64 {
	return meters * 100
}

func centimetersToDecimeters(centimeters float64) float64 {
	return centimeters / 10
}

func decimetersToCentimeters(decimeters float64) float64 {
	return decimeters * 10
}

func metersToDecimeters(meters float64) float64 {
	return meters * 10
}

func decimetersToMeters(decimeters float64) float64 {
	return decimeters / 10
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	from := queryParams.Get("from")
	to := queryParams.Get("to")
	valueStr := queryParams.Get("value")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var result float64
	var fromUnit string
	var toUnit string

	switch from {
	case "centimeters":
		fromUnit = "centimeters"
		switch to {
		case "meters":
			toUnit = "meters"
			result = centimetersToMeters(value)
		case "decimeters":
			toUnit = "decimeters"
			result = centimetersToDecimeters(value)
		default:
			http.Error(w, "Invalid conversion", http.StatusBadRequest)
			return
		}
	case "decimeters":
		fromUnit = "decimeters"
		switch to {
		case "meters":
			toUnit = "meters"
			result = decimetersToMeters(value)
		case "centimeters":
			toUnit = "centimeters"
			result = decimetersToCentimeters(value)
		default:
			http.Error(w, "Invalid conversion", http.StatusBadRequest)
			return
		}
	case "meters":
		fromUnit = "meters"
		switch to {
		case "centimeters":
			toUnit = "centimeters"
			result = metersToCentimeters(value)
		case "decimeters":
			toUnit = "decimeters"
			result = metersToDecimeters(value)
		default:
			http.Error(w, "Invalid conversion", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Invalid conversion", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "%.2f %s is equal to %.2f %s", value, fromUnit, result, toUnit)
}

func main() {
	http.HandleFunc("/convert", convertHandler)
	http.ListenAndServe(":8000", nil)
}
