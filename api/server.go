package api

import (
	"encoding/json"
	"log"
	"net/http"
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

	mux.HandleFunc("/", s.handleNotFound)

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
		"time": time.Now().Format(time.RFC3339),
		"battery": battery,
		"app_name": "freeport",
	}

	w.Header().Set("Content-Type", "applications/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}