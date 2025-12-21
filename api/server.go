package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/system/battery", s.handleBattery)

	mux.HandleFunc("/", s.handleCustomOrNotFound)

	log.Printf("API Server starting on port %s", s.port)
	return http.ListenAndServe(":"+s.port, s.methodFilter(mux))
}

func (s *Server) methodFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleBattery(w http.ResponseWriter, r *http.Request) {
	battery, err := GetBatteryPercentage()
	if err != nil {
		http.Error(w, "Failed to get battery info", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"time":     time.Now().Format(time.RFC3339),
		"battery":  battery,
		"app_name": "freeport",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCustomOrNotFound(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if path == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if len(parts) >= 2 {
		appName := parts[0]
		endpoint := parts[1]

		if endpoint == "init" {
			s.handleCustomInit(w, r, appName)
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}