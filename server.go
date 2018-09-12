package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

//Server Server for allowing application to process HTTP requests
type Server struct {
	config *HTTPConfig
	mux    *http.ServeMux
	app    *Application
}

//InitializeServer Creates a server with all attached dependencies
func InitializeServer(config *HTTPConfig, app *Application) *Server {
	mux := http.NewServeMux()
	server := &Server{config, mux, app}
	server.initRoutes()
	return server
}

//Run Runs the Server by listening on the configured port
func (s *Server) Run() {
	portString := strconv.Itoa(s.config.Port)
	log.Println("Listening for traffic on port " + portString)
	http.ListenAndServe(":"+portString, s.mux)
}

func (s *Server) initRoutes() {
	s.mux.HandleFunc("/upload", logRequest(validateUpload(s.handleUpload)))
}

//UploadResponse The json response returned to the Client
type UploadResponse struct {
	Art   []string
	Error []string
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		duration := time.Since(start)
		durationString := fmt.Sprintf("%0.3f", duration.Seconds())
		log.Println(r.Method + " request received from " + r.RemoteAddr + ", took " + durationString + " seconds")
	})
}

func validateUpload(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, _, err := r.FormFile("image")

		if err != nil {
			response := UploadResponse{Art: []string{}, Error: []string{err.Error()}}
			jsonResponse, _ := json.MarshalIndent(response, "", "  ")
			w.WriteHeader(400)
			w.Write([]byte(jsonResponse))
			return
		}
		next(w, r)

	})
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	data, _, err := r.FormFile("image")

	lines, err := s.app.ProcessData(data)

	var response UploadResponse
	if err == nil {
		response = UploadResponse{Art: lines, Error: []string{}}
		w.WriteHeader(200)
	} else {
		response = UploadResponse{Art: []string{}, Error: []string{err.Error()}}
		w.WriteHeader(500)
	}
	jsonResponse, _ := json.MarshalIndent(response, "", "  ")
	w.Write([]byte(jsonResponse))

}
