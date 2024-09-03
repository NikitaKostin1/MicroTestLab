package main

import (
	"net/http"
)



func (app *AppConfig) HandleBrokerRequest(writer http.ResponseWriter, request *http.Request) {
	payload := APIResponse{
		IsError: false,
		Message: "Broker request received successfully",
	}

	err := app.writeJSON(writer, http.StatusOK, payload)
	if err != nil {
		http.Error(writer, "Unable to write JSON response", http.StatusInternalServerError)
	}
}
