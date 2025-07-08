package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"juggler/internal/juggler"
	"juggler/internal/web"
)

func TestWebServerStats(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleStats)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	expected := "application/json"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("Expected content type %s, got %s", expected, ct)
	}

	var stats web.StatsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if stats.InHand != 0 {
		t.Errorf("Expected 0 balls in hand, got %d", stats.InHand)
	}

	if stats.InAir != 0 {
		t.Errorf("Expected 0 balls in air, got %d", stats.InAir)
	}

	if stats.TotalBalls != 0 {
		t.Errorf("Expected 0 total balls, got %d", stats.TotalBalls)
	}

	if stats.IsRunning {
		t.Error("Expected juggling to not be running")
	}
}

func TestWebServerStart(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	startReq := web.StartRequest{
		TotalBalls:  3,
		TimeMinutes: 2,
	}

	jsonBody, _ := json.Marshal(startReq)
	req, err := http.NewRequest("POST", "/api/start", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleStart)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if response["status"] != "started" {
		t.Errorf("Expected status 'started', got %s", response["status"])
	}

	if j.GetTotalBalls() != 3 {
		t.Errorf("Expected 3 balls after start, got %d", j.GetTotalBalls())
	}
}

func TestWebServerStartInvalidMethod(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	req, err := http.NewRequest("GET", "/api/start", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleStart)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestWebServerStartInvalidJSON(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	req, err := http.NewRequest("POST", "/api/start", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleStart)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, status)
	}
}

func TestWebServerStartInvalidValues(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	tests := []struct {
		name     string
		balls    int
		time     int
		expected int
	}{
		{"Zero balls", 0, 2, http.StatusBadRequest},
		{"Negative balls", -1, 2, http.StatusBadRequest},
		{"Zero time", 3, 0, http.StatusBadRequest},
		{"Negative time", 3, -1, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startReq := web.StartRequest{
				TotalBalls:  tt.balls,
				TimeMinutes: tt.time,
			}

			jsonBody, _ := json.Marshal(startReq)
			req, err := http.NewRequest("POST", "/api/start", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(server.HandleStart)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expected {
				t.Errorf("Expected status code %d, got %d", tt.expected, status)
			}
		})
	}
}

func TestWebServerStop(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	j.Reset(3, 2)

	req, err := http.NewRequest("POST", "/api/stop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleStop)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if response["status"] != "stopped" {
		t.Errorf("Expected status 'stopped', got %s", response["status"])
	}

	if !j.IsFinished() {
		t.Error("Expected juggler to be finished after stop")
	}
}

func TestWebServerStopInvalidMethod(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	req, err := http.NewRequest("GET", "/api/stop", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleStop)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestWebServerHome(t *testing.T) {
	j := juggler.NewJuggler(0, 0)
	server := web.NewServer(j, 8080)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleHome)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	expected := "text/html; charset=utf-8"
	if ct := rr.Header().Get("Content-Type"); ct != expected {
		t.Errorf("Expected content type %s, got %s", expected, ct)
	}

	body := rr.Body.String()
	if !containsHTML(body) {
		t.Error("Expected response to contain HTML")
	}
}

func TestStatsResponseStructure(t *testing.T) {
	j := juggler.NewJuggler(3, 2)
	j.Reset(3, 2)

	server := web.NewServer(j, 8080)

	req, err := http.NewRequest("GET", "/api/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.HandleStats)
	handler.ServeHTTP(rr, req)

	var stats web.StatsResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Errorf("Failed to parse JSON response: %v", err)
	}

	if stats.InHand != 3 {
		t.Errorf("Expected 3 balls in hand, got %d", stats.InHand)
	}

	if stats.InAir != 0 {
		t.Errorf("Expected 0 balls in air, got %d", stats.InAir)
	}

	if stats.TotalBalls != 3 {
		t.Errorf("Expected 3 total balls, got %d", stats.TotalBalls)
	}

	if stats.TotalTime != 2 {
		t.Errorf("Expected 2 minutes total time, got %d", stats.TotalTime)
	}

	if len(stats.Balls) != 3 {
		t.Errorf("Expected 3 balls in details, got %d", len(stats.Balls))
	}

	for i, ball := range stats.Balls {
		if ball.ID <= 0 {
			t.Errorf("Ball %d has invalid ID: %d", i, ball.ID)
		}
		if ball.Status != "in_hand" {
			t.Errorf("Ball %d expected to be in hand, got status: %s", i, ball.Status)
		}
	}
}

func containsHTML(s string) bool {
	return bytes.Contains([]byte(s), []byte("<html>")) &&
		bytes.Contains([]byte(s), []byte("</html>")) &&
		bytes.Contains([]byte(s), []byte("<body>")) &&
		bytes.Contains([]byte(s), []byte("</body>"))
}
