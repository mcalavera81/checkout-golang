package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Details string `json:details`
	StatusCode int `json:"status"`
}

func writeError(w io.Writer, err error, message string, statusCode int){
	error := Error{Message: message, Details: err.Error(), StatusCode: statusCode, }
	json.NewEncoder(w).Encode(&error)
}

func errorHandler(w http.ResponseWriter, err error, status int, message string){
	w.WriteHeader(status)
	writeError(w, err, message, status)
}