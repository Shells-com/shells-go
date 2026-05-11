package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"sync"
)

// ClaudeAPIServer provides an HTTP API for Claude to control the remote desktop
type ClaudeAPIServer struct {
	claudeCtrl *ClaudeControlInterface
	server     *http.Server
	mutex      sync.Mutex
}

// NewClaudeAPIServer creates a new API server for Claude Computer Use
func NewClaudeAPIServer(claudeCtrl *ClaudeControlInterface, port int) *ClaudeAPIServer {
	api := &ClaudeAPIServer{
		claudeCtrl: claudeCtrl,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/control", api.handleControl)
	mux.HandleFunc("/api/v1/screenshot", api.handleScreenshot)
	mux.HandleFunc("/api/v1/status", api.handleStatus)

	api.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}

	return api
}

// Start launches the API server
func (api *ClaudeAPIServer) Start() error {
	log.Printf("Starting Claude API server on %s", api.server.Addr)
	return api.server.ListenAndServe()
}

// Stop shuts down the API server
func (api *ClaudeAPIServer) Stop() error {
	return api.server.Close()
}

// handleControl processes control API requests
func (api *ClaudeAPIServer) handleControl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req APIRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request format: %v", err), http.StatusBadRequest)
		return
	}

	result, err := api.claudeCtrl.HandleAPIRequest(string(mustMarshal(req)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing request: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

// handleScreenshot returns the current screenshot as a base64-encoded JPEG
func (api *ClaudeAPIServer) handleScreenshot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	screenshotImg := api.claudeCtrl.GetScreenshot()
	if screenshotImg == nil {
		http.Error(w, "No screenshot available", http.StatusNotFound)
		return
	}

	// Encode image as JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, screenshotImg, &jpeg.Options{Quality: 75}); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding image: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode as base64
	base64Data := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Return the data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"image": base64Data,
		"format": "image/jpeg;base64",
	})
}

// handleStatus returns the current status of the Claude Computer Use integration
func (api *ClaudeAPIServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := map[string]interface{}{
		"enabled": api.claudeCtrl.enabled,
		"connected": api.claudeCtrl.inputs != nil,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Helper function to marshal JSON
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}