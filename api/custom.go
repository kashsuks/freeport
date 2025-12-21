package api

import (
	"encoding/json"
	"net/http"
	"sync"
)

type CustomProtocol struct {
	AppName string
	Passkey string
	Description string
}

var (
	protocols = make(map[string]CustomProtocol)
	mu sync.RWMutex
)

func RegisterProtocol(appName, passkey, description string) {
	mu.Lock()
	defer mu.Unlock()
	protocols[appName] = CustomProtocol{
		AppName: appName,
		Passkey: passkey,
		Description: description,
	}
}

func ValidateProtocol(appName, passkey string) bool {
	mu.RLock()
	defer mu.RUnlock()
	protocol, exists := protocols[appName]
	if !exists {
		return false
	}
	return protocol.Passkey == passkey
}

func (s *Server) handleCustomInit(w http.ResponseWriter, r *http.Request, appName string) {
	headerAppName := r.Header.Get("X-App-Name")
	headerPasskey := r.Header.Get("X-Passkey")

	if headerAppName != appName {
		http.Error(w, "App name mismatch", http.StatusBadRequest)
		return
	}

	if !ValidateProtocol(appName, headerPasskey) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"message": "Hello, World!",
		"app_name": appName,
		"status": "initialized",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}