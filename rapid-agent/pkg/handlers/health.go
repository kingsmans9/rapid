package handlers

import (
	"net/http"
)

type HealthzResponse struct {
	Status string `json:"status"`
}

func Healthz(w http.ResponseWriter, r *http.Request) {

	// TODO: Check if the database is ready
	healthzResponse := HealthzResponse{
		Status: "OK",
	}

	JSON(w, http.StatusOK, healthzResponse)
}
