// main.go
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type conversionFunc func(float64) float64

var unitConversions = map[string]map[string]conversionFunc{
	"centimeters": {
		"meters":      func(cm float64) float64 { return cm / 100 },
		"decimeters":  func(cm float64) float64 { return cm / 10 },
	},
	"decimeters": {
		"meters":      func(dm float64) float64 { return dm / 10 },
		"centimeters": func(dm float64) float64 { return dm * 10 },
	},
	"meters": {
		"centimeters": func(m float64) float64 { return m * 100 },
		"decimeters":  func(m float64) float64 { return m * 10 },
	},
}

type PageVariables struct {
	From   string
	To     string
	Value  float64
	Result float64
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	from, to, valueStr := queryParams.Get("from"), queryParams.Get("to"), queryParams.Get("value")

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	conversion, ok := unitConversions[from][to]
	if !ok {
		http.Error(w, "Invalid conversion", http.StatusBadRequest)
		return
	}

	result := conversion(value)

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := PageVariables{
		From:   from,
		To:     to,
		Value:  value,
		Result: result,
	}

	tmpl.Execute(w, data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func main() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
		<-ch
		fmt.Println("Received interrupt signal. Shutting down...")
		os.Exit(0)
	}()

	http.HandleFunc("/convert", convertHandler)
	http.HandleFunc("/", indexHandler)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
