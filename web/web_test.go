package web

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonResponse(t *testing.T) {
	// Initialize test data
	code := 200
	data := "some response data"
	errMsg := "some error message"
	w := httptest.NewRecorder()

	// Call the function
	jsonResponse(code, data, nil, w)

	// Check response
	if w.Code != code {
		t.Errorf("Expected status code %d, but got %d", code, w.Code)
	}
	expectedBody := `{"Data":"some response data","Error":""}
`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
	}

	// Call the function with error
	w = httptest.NewRecorder()
	jsonResponse(code, data, fmt.Errorf(errMsg), w)

	// Check response
	if w.Code != code {
		t.Errorf("Expected status code %d, but got %d", code, w.Code)
	}
	expectedBody = `{"Data":"some response data","Error":"some error message"}
`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
	}
}

func TestStart(t *testing.T) {
	// Initialize test data
	port := "8080"
	handlers := http.HandlerFunc(okHandler)

	// Call the function
	//var stopCh chan<- struct{}
	var err error
	go func() {
		_, err = Start(handlers, port, nil)
	}()

	// Check error
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Send a request to the server and check the response
	req, err := http.NewRequest(http.MethodGet, "http://localhost:"+port+"/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}
	expectedBody := `{"Data":"Hello","Error":""}
`
	if body, _ := io.ReadAll(resp.Body); string(body) != expectedBody {
		t.Errorf("Expected body %s, but got %s", expectedBody, string(body))
	}

	// Stop the server and check the stop channel
	//	stopCh <- struct{}{}
}

func TestRegisterPath(t *testing.T) {
	// Initialize test data
	path := "/test"
	url := "/path/to/files"
	mux := http.NewServeMux()

	// Call the function
	result := registerPath(path, url, nil, mux)

	// Check result
	if result == nil {
		t.Errorf("Expected object, but got nil")
	}
	if _, pattern := result.Handler(httptest.NewRequest("GET", "http://localhost/test", nil)); pattern != "/test/" {
		t.Errorf("Expected pattern %s, but got %s", "/test/", pattern)
	}
}

func TestOkHandler(t *testing.T) {
	// Initialize test data
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Call the function
	okHandler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}
	expectedBody := `{"Data":"Hello","Error":""}
`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, but got %s", expectedBody, w.Body.String())
	}
}
