package server

import (
	"encoding/json"
	"net/http"
)

func WriteMessage(wr http.ResponseWriter, code int, message string) {
	var data []byte
	if code >= 200 && code < 300 {
		var msg struct {
			Message string `json:"message"`
		}
		msg.Message = message
		data, _ = json.Marshal(msg)
	} else if code >= 400 && code < 600 {
		var resp struct {
			Error string `json:"error"`
		}
		resp.Error = message
		data, _ = json.Marshal(resp)
	} else {
		return
	}
	WriteResponse(wr, code, data)
}

func WriteResponse(wr http.ResponseWriter, code int, data []byte) {
	wr.Header().Set("Content-Type", "application/json")
	wr.WriteHeader(code)
	wr.Write(data)
}
