package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorBody struct {
	Error string `json:"error"`
}

func Failed(writer http.ResponseWriter, status int, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	e := new(errorBody)
	e.Error = err.Error()

	resp, err := json.Marshal(&e)
	if err != nil {
		go log.Println("failed to parse error message: " + err.Error())
		return
	}
	if _, err := writer.Write(resp); err != nil {
		go log.Println("failed to send error message: " + err.Error())
	}
}

func Success(writer http.ResponseWriter, status int, resp interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	bytesResp, err := json.Marshal(resp)
	if err != nil {
		Failed(writer, http.StatusInternalServerError, err)
		return
	}
	if _, err := writer.Write(bytesResp); err != nil {
		log.Println("failed to send success message: " + err.Error())
	}
}
