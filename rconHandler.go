package main

import (
	"encoding/json"
	"net/http"

	rcon "github.com/gtaylor/factorio-rcon"
)

// RCONHandler encapsulates the rcon.RCON client for handling commands
type RCONHandler struct {
	client *rcon.RCON
}

type commandPayload struct {
	Command string
}

type commandResponse struct {
	Result string `json:"result"`
}

func (rh *RCONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var requestPayload commandPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, "Cannot parse body as JSON", 400)
		return
	}
	defer r.Body.Close()

	response, err := rh.client.Execute(requestPayload.Command)
	if err != nil {
		http.Error(w, "Cannot execute command", 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(commandResponse{
		Result: response.Body,
	})
}
