package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEchoHandler(t *testing.T) {
	// Create a new request with a JSON-encoded body
	requestBody := map[string]string{"message": "hello world"}
	requestBodyBytes, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/echo", bytes.NewBuffer(requestBodyBytes))

	// Create a new response recorder to capture the handler's output
	rr := httptest.NewRecorder()

	// Call the echoHandler function with the request and response recorder
	echoHandler(rr, req)

	// Check the response status code is 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body matches the input body
	// Leave the newline!
	expectedResponseBody := `{"data":{"echo":"{\"message\":\"hello world\"}"}}
`
	if responseBody := rr.Body.String(); responseBody != expectedResponseBody {
		t.Errorf("handler returned unexpected body: got %v want %v", responseBody, expectedResponseBody)
	}
}
