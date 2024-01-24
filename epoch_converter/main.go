package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type ConversionResult struct {
	Result string `json:"result"`
}

func main() {
	http.HandleFunc("/convert", ConvertHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))

	port := ":8080"
	println("Server started on port", port)
	http.ListenAndServe(port, nil)
}

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
	timestampStr := r.URL.Query().Get("timestamp")
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid timestamp", http.StatusBadRequest)
		return
	}

	convertedTime := time.Unix(timestamp, 0)
	gmtTime := convertedTime.UTC()

	localTime, _ := time.LoadLocation("Asia/Kolkata") // Adjust the time zone as needed
	localTimeStr := gmtTime.In(localTime).Format("Monday, 02 January 2006 15:04:05 MST")

	// Calculate relative time
	relativeTime := time.Since(convertedTime)
	relativeTimeString := humanizeDuration(relativeTime)

	// Format the result as a JSON object
	result := ConversionResult{
		Result: fmt.Sprintf("GMT: %s\nYour time zone: %s\nRelative: %s", gmtTime.Format("Monday, 02 January 2006 15:04:05 MST"), localTimeStr, relativeTimeString),
	}

	// Convert the result to JSON
	jsonResult, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

// Function to humanize duration in a user-friendly format
func humanizeDuration(duration time.Duration) string {
	if duration < time.Minute {
		return "A few seconds ago"
	} else if duration < 2*time.Minute {
		return "A minute ago"
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	} else if duration < 2*time.Hour {
		return "An hour ago"
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	}

	return ""
}
