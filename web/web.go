package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	r := mux.NewRouter()

	// Echo endpoint
	r.HandleFunc("/echo", echoHandler).Methods("GET")

	// Start server
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", r)
}

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		jsonResponseHandler(w, nil, http.StatusBadRequest, err)
		return
	}

	responseData := map[string]string{"echo": string(data)}
	jsonResponseHandler(w, responseData, http.StatusOK, nil)
}

func jsonResponseHandler(w http.ResponseWriter, data interface{}, statusCode int, err error) {
	var response *Response
	if err != nil {
		response = &Response{Error: err.Error()}
	} else {
		response = &Response{Data: data}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
