package web

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Data  interface{}
	Error string
}

// jsonResponse takes a response code, message interface, error, and ResponseWriter object and writes the data
// and response code to the ResponseWriter as JSON using a Response struct that contains optional
// data and error fields
func jsonResponse(code int, message interface{}, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := Response{
		Data:  message,
		Error: "",
	}

	if err != nil {
		response.Error = err.Error()
	}

	json.NewEncoder(w).Encode(response)
}

// Start creates and starts a web server that routes requests to http handlers. The http handlers are passed in as a parameter as
// well as the port and TLS configuration for the server. It returns an error and a stop channel to stop the server.
func Start(handlers http.Handler, port string, tlsConfig *tls.Config) (stopChan chan<- struct{}, err error) {
	mux := http.NewServeMux()
	mux.Handle("/", handlers)

	server := &http.Server{
		Addr:      ":" + port,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	stopCh := make(chan struct{})
	go func() {
		<-stopCh

		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown failed: %v", err)
		}
	}()

	if tlsConfig == nil {
		if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return nil, err
		}
	} else {
		if err = server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			return nil, err
		}
	}

	return stopCh, nil
}

// registerPath is a function that takes a path, URL and http handelr and adds it to a muxer, it then returns the
// muxer as an object
func registerPath(path string, url string, handler http.Handler, mux *http.ServeMux) *http.ServeMux {
	if path != "/" && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	mux.Handle(path+"/", http.StripPrefix(path, http.FileServer(http.Dir(url))))
	return mux
}

// okHandler is an http Handler that just returns 200 OK and a "Hello" response when accessed.
func okHandler(w http.ResponseWriter, r *http.Request) {
	jsonResponse(http.StatusOK, "Hello", nil, w)
}
