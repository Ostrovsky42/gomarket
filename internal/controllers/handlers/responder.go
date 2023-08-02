package handlers

import (
	"encoding/json"
	"gomarket/internal/logger"
	"net/http"
)

const (
	ContentType = "Content-Type"
	JSON        = "application/json"
)

func writeOKWithJSON(w http.ResponseWriter, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to marshal JSON response")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set(ContentType, JSON)
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(jsonData); err != nil {
		logger.Log.Error().Err(err).Msg("failed write JSON data")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
