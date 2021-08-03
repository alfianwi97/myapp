package index

import (
	"net/http"

	"github.com/alfianwi97/myapp/pkg/router"
)

// GetIndex Function to Show API Information
func GetIndex(w http.ResponseWriter, r *http.Request) {
	router.ResponseSuccess(w, "Codebase Go REST is running")
}

// GetHealth Function to Show Health Check Status
func GetHealth(w http.ResponseWriter, r *http.Request) {
	router.HealthCheck(w)
}
