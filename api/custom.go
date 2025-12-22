package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type CustomProtocol struct {
	AppName string
	Passkey string
	Description string
	Methods map[string]string
	Data map[string]interface{}
	History map[string][]DataEntry
}

type DataEntry struct {
	Data interface{}
	Timestamp time.Time
	Source string
}

var (
	protocols = make(map[string]*CustomProtocol)
	mu sync.RWMutex
)

func RegisterProtocol(appName, passkey, description string) {
	mu.Lock()
	defer mu.Unlock()
	protocols[appName] = &CustomProtocol{
		AppName: appName,
		Passkey: passkey,
		Description: description,
		Methods: make(map[string]string),
		Data: make(map[string]interface{}),
		History: make(map[string][]DataEntry),
	}
	protocols[appName].Methods["init"] = "Initialize connection"
}

func RegisterMethod(appName, methodName, description string) {
	mu.Lock()
	defer mu.Unlock()
	if protocol, exists := protocols[appName]; exists {
		protocol.Methods[methodName] = description
	}
}

func StoreData(appName, methodName, source string, data interface{}) bool {
	mu.Lock()
	defer mu.Unlock()
	if protocol, exists := protocols[appName]; exists {
		protocol.Data[methodName] = data

		entry := DataEntry{
			Data: data,
			Timestamp: time.Now(),
			Source: source,
		}
		protocol.History[methodName] = append(protocol.History[methodName], entry)

		if len(protocol.History[methodName]) > 100 {
			protocol.History[methodName] = protocol.History[methodName][1:]
		}

		return true
	}
	return false
}

func GetData(appName, methodName string) (interface{}, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if protocol, exists := protocols[appName]; exists {
		data, ok := protocol.Data[methodName]
		return data, ok
	}
	return nil, false
}

func GetHistory(appName, methodName string, limit int) ([]DataEntry, bool) {
	mu.RLock()
	defer mu.RUnlock()
	if protocol, exists := protocols[appName]; exists {
		history, ok := protocol.History[methodName]
		if !ok {
			return nil, false
		}

		start := 0
		if len(history) > limit {
			start = len(history) - limit
		}
		return history[start:], true
	}
	return nil, false
}

func ClearData(appName, methodName string) bool {
	mu.Lock()
	defer mu.Unlock()
	if protocol, exists := protocols[appName]; exists {
		delete(protocol.Data, methodName)
		delete(protocol.History, methodName)
		return true
	}
	return false
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

func MethodExists(appName, methodName string) bool {
	mu.RLock()
	defer mu.RUnlock()
	if protocol, exists := protocols[appName]; exists {
		_, methodExists := protocol.Methods[methodName]
		return methodExists
	}
	return false
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
		"time": time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleCustomMethod(w http.ResponseWriter, r *http.Request, appName, methodName string) {
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

	if !MethodExists(appName, methodName) {
		http.Error(w, "Method not found", http.StatusNotFound)
		return
	}

	if r.Method == http.MethodPost {
		var requestData map[string]interface{}
		
		if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		source := headerAppName
		if src, ok := requestData["source"]; ok {
			source = fmt.Sprintf("%v", src)
		}

		if StoreData(appName, methodName, source, requestData) {
			response := map[string]interface{}{
				"status":    "success",
				"message":   "Data stored successfully",
				"app_name":  appName,
				"method":    methodName,
				"timestamp": time.Now().Format(time.RFC3339),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Failed to store data", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodGet {
		data, exists := GetData(appName, methodName)
		
		response := map[string]interface{}{
			"app_name": appName,
			"method":   methodName,
			"time":     time.Now().Format(time.RFC3339),
		}

		if exists {
			response["data"] = data
			response["status"] = "success"
		} else {
			response["message"] = "No data available"
			response["status"] = "no_data"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if r.Method == http.MethodDelete {
		if ClearData(appName, methodName) {
			response := map[string]interface{}{
				"status":   "success",
				"message":  "Data cleared successfully",
				"app_name": appName,
				"method":   methodName,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, "Failed to clear data", http.StatusInternalServerError)
		}
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (s *Server) handleCustomHistory(w http.ResponseWriter, r *http.Request, appName, methodName string) {
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

	history, exists := GetHistory(appName, methodName, 10)

	if !exists {
		http.Error(w, "No history available", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"app_name": appName,
		"method": methodName,
		"count": len(history),
		"history": history,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}