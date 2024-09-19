package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)



type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	IsError bool        `json:"is_error"`
}



func (app *AppConfig) readJSON(writer http.ResponseWriter, request *http.Request, data any) error {
	const maxRequestBodySize = 1 << 20 // 1 MB

	request.Body = http.MaxBytesReader(writer, request.Body, maxRequestBodySize)

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(data)
	if err != nil {
		return err
	}

	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("The body should contain only one JSON value.")
	}

	return nil
}

func (app *AppConfig) writeJSON(writer http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			writer.Header()[key] = value
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (app *AppConfig) errorJSON(writer http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	responsePayload := APIResponse{
		IsError:   true,
		Message:   err.Error(),
	}

	return app.writeJSON(writer, statusCode, responsePayload)
}
