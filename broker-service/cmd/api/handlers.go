package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"fmt"
)



type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}



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

func (app *AppConfig) HandleSubmission(writer http.ResponseWriter, request *http.Request) {
	var reqPayload RequestPayload

	err := app.readJSON(writer, request, &reqPayload)
	if err != nil {
		app.errorJSON(writer, err, http.StatusBadRequest)
		return
	}

	switch reqPayload.Action {
	case "auth":
		app.handleAuth(writer, reqPayload.Auth)
	default:
		app.errorJSON(writer, errors.New("Unknown action in request"), http.StatusBadRequest)
	}
}

func (app *AppConfig) handleAuth(writer http.ResponseWriter, payload AuthPayload) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		app.errorJSON(writer, err, http.StatusInternalServerError)
		return
	}

	url := "http://authentication-service:1025/authenticate"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(writer, err, http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	clientResponse, err := client.Do(request)
	if err != nil {
		app.errorJSON(writer, err, http.StatusInternalServerError)
		return
	}
	defer clientResponse.Body.Close()

	if clientResponse.StatusCode == http.StatusUnauthorized {
		app.errorJSON(writer, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	} else if clientResponse.StatusCode != http.StatusAccepted {
		app.errorJSON(writer, fmt.Errorf("unexpected status code: %d", clientResponse.StatusCode), http.StatusInternalServerError)
		return
	}

	var authResponse APIResponse
	err = json.NewDecoder(clientResponse.Body).Decode(&authResponse)
	if err != nil {
		app.errorJSON(writer, err, http.StatusInternalServerError)
		return
	}

	if authResponse.IsError {
		app.errorJSON(writer, errors.New(authResponse.Message), http.StatusUnauthorized)
		return
	}

	response := APIResponse{
		Message: "User authenticated successfully",
		Data:    authResponse.Data,
		IsError: false,
	}

	app.writeJSON(writer, http.StatusAccepted, response)
}
