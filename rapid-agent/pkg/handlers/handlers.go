package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pavansokkenagaraj/rapid-agent/pkg/util/logger"
)

type ResponseErr struct {
	Message string `json:"message,omitempty"`
}

func JSONWithError(w http.ResponseWriter, status int, errMsg string) {
	logger.Errorf(errMsg)
	var responseErr ResponseErr
	if status != http.StatusInternalServerError {
		responseErr = ResponseErr{Message: errMsg}
	}

	JSON(w, status, responseErr)
}

func JSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}
