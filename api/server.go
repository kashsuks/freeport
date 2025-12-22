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
	return http.ListenAndServe(":"+s.port, mux)
}

func (s *Server) handleBattery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if len(parts) < 2 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	appName := parts[0]
	methodName := parts[1]

	if appName == "system" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}

	if methodName == "init" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed for init", http.StatusMethodNotAllowed)
			return
		}
		s.handleCustomInit(w, r, appName)
		return
	}

	if len(parts) == 3 && parts[2] == "history" {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		s.handleCustomHistory(w, r, appName, methodName)
		return
	}

	if len(parts) == 2 {
		if r.Method == http.MethodGet || r.Method == http.MethodPost || r.Method == http.MethodDelete {
			s.handleCustomMethod(w, r, appName, methodName)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}