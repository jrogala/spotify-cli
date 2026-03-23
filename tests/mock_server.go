// Package tests provides a mock Spotify API server for integration tests.
package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
)

// RecordedRequest stores details of a received HTTP request.
type RecordedRequest struct {
	Method string
	Path   string
	Body   string
}

// MockServer simulates the Spotify API for testing.
type MockServer struct {
	Server   *httptest.Server
	mu       sync.Mutex
	routes   map[string]mockRoute
	Requests []RecordedRequest
}

type mockRoute struct {
	statusCode int
	body       string
}

func routeKey(method, path string) string {
	return method + " " + path
}

// NewMockServer creates and starts a mock HTTP server.
func NewMockServer() *MockServer {
	ms := &MockServer{
		routes: make(map[string]mockRoute),
	}

	ms.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ms.mu.Lock()
		ms.Requests = append(ms.Requests, RecordedRequest{
			Method: r.Method,
			Path:   r.URL.RequestURI(),
		})

		// Try exact match first (with query string)
		key := routeKey(r.Method, r.URL.RequestURI())
		route, ok := ms.routes[key]
		if !ok {
			// Try path-only match (without query string)
			key = routeKey(r.Method, r.URL.Path)
			route, ok = ms.routes[key]
		}
		ms.mu.Unlock()

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"status":  404,
					"message": "mock: no route for " + r.Method + " " + r.URL.RequestURI(),
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(route.statusCode)
		w.Write([]byte(route.body))
	}))

	return ms
}

// On registers a response for a given method and path.
func (ms *MockServer) On(method, path string, statusCode int, responseBody string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.routes[routeKey(method, path)] = mockRoute{
		statusCode: statusCode,
		body:       responseBody,
	}
}

// OnJSON registers a response, marshalling body to JSON.
func (ms *MockServer) OnJSON(method, path string, statusCode int, body any) {
	data, _ := json.Marshal(body)
	ms.On(method, path, statusCode, string(data))
}

// Reset clears all registered routes and recorded requests.
func (ms *MockServer) Reset() {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.routes = make(map[string]mockRoute)
	ms.Requests = nil
}

// Close shuts down the server.
func (ms *MockServer) Close() {
	ms.Server.Close()
}

// BaseURL returns the server's base URL.
func (ms *MockServer) BaseURL() string {
	return ms.Server.URL
}
