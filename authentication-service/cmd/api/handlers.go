package main

import (
	"errors"
	"fmt"
	"net/http"
)



func (app *AppConfig) HandleAuth(writer http.ResponseWriter, request *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(writer, request, &requestPayload)
	if err != nil {
		app.errorJSON(writer, err, http.StatusBadRequest)
		return
	}

	user, err := app.DBModels.User.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(writer, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	isValid, err := user.IsPasswordMatching(requestPayload.Password)
	if err != nil || !isValid {
		app.errorJSON(writer, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	responsePayload := APIResponse{
		IsError: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(writer, http.StatusAccepted, responsePayload)
}
