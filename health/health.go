// Package health provides functionality for constructing a health check endpoint.
package health

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
	// OK is a healthy status
	OK = "Healthy"
	// WARNING is a warning status
	WARNING = "Warning"
	// CRITICAL is a critical status
	CRITICAL = "Critical"
)

// Response represents the top-level data structure returned from the API
type Response struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Version string `json:"version"`
	Checks  Checks `json:"checks"`
}

// Check represents each individual health check
type Check struct {
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

// Checks are a set of named Check pointers
type Checks []*Check

// Checker is an interface for returning healthchecks
type Checker interface {
	GetChecks() Checks
	GetResponse(Checks) *Response
}

// AddServiceCheck inspects all other checks in the map and adds a "service" entry with the highest
// error level found.
func AddServiceCheck(checks Checks) {
	for _, c := range checks {
		if c.Name == "service" {
			return
		}
	}

	service := &Check{
		Name:        "service",
		Status:      OK,
		Description: "Service OK",
	}

	// If any value is not OK, upgrade the value of "service" to match.
	for _, check := range checks {
		if check == nil {
			continue
		}

		if check.Status == CRITICAL || check.Status == WARNING && service.Status != CRITICAL {
			service = check
		}
	}

	checks = append(checks, service)
}

// CreateHealthHandler takes a health checker and returns an endpoint for health checks
func CreateHealthHandler(checker Checker) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {

		writer.Header().Set("Cache-Control", "no-store")

		res := checker.GetResponse(checker.GetChecks())

		if res.Status == "" {
			res.Status = findStatus(res.Checks)
		}

		serveJSON(writer, reader, http.StatusOK, res)
	}
}

// CreateProbeHandler takes a health checker and returns an endpoint for probe
// checks which return different http statuses
func CreateProbeHandler(checker Checker) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, reader *http.Request) {

		writer.Header().Set("Cache-Control", "no-store")

		checks := checker.GetChecks()

		res := checker.GetResponse(checks)

		// Return OK by default
		status := http.StatusOK

		// Determine if we should just return http.StatusOK or http.StatusServiceUnavailable
		// iterate over checker.GetChecks
		for _, value := range checks {
			// If any check is not OK, status will be set to http.StatusServiceUnavailable
			if value.Status != OK {
				status = http.StatusServiceUnavailable
			}
		}

		if res.Status == "" {
			res.Status = findStatus(res.Checks)
		}

		serveJSON(writer, reader, status, res)
	}
}

func findStatus(checks Checks) string {
	for _, value := range checks {
		// If any check is not OK, status will be set to http.StatusServiceUnavailable
		if value.Status == WARNING {
			return WARNING
		}

		if value.Status == CRITICAL {
			return CRITICAL
		}
	}

	return OK
}

func serveJSON(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(data)
	if err != nil {
		log.Printf("error: could not encode JSON: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(statusCode)
	if _, err := io.WriteString(w, string(j)); err != nil {
		log.Printf("error: could not write to ResponseWriter: %s\n", err)
		return
	}
}
