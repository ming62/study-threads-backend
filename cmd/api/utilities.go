package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
    wrapper := make(map[string]interface{})

    wrapper[wrap] = data

    js, err := json.Marshal(wrapper)
    if err != nil {
        return err
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _, err = w.Write(js)
    if err != nil {
        return err
    }

    return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {

	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}
	
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	app.writeJSON(w, statusCode, theError, "error")
}
