package eureka

import (
	"encoding/base64"
	"encoding/json"
	_ "fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterWithEureka(t *testing.T) {
	// Mock Eureka server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request has the correct authorization header
		authHeader := r.Header.Get("Authorization")
		expectedAuthHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte("testuser:testpassword"))

		if authHeader != expectedAuthHeader {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the request body is correct
		var requestData map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		expectedInstanceID := "test-instance"
		if requestData["instance"].(map[string]interface{})["instanceId"].(string) != expectedInstanceID {
			http.Error(w, "Invalid instance ID in request body", http.StatusBadRequest)
			return
		}

		// Send a successful response
		w.WriteHeader(http.StatusNoContent)
	}))

	defer mockServer.Close()

	// Create a test configuration
	config := Config{
		EurekaURL: mockServer.URL,
		App:       "test-app",
		Instance: InstanceData{
			InstanceID: "test-instance",
		},
	}

	// Perform the test
	err := RegisterWithEureka(config, "testuser", "testpassword")
	if err != nil {
		t.Errorf("RegisterWithEureka failed: %v", err)
	}
}

func TestRegisterWithEurekaUnauthorized(t *testing.T) {
	// Mock Eureka server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send an unauthorized response
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))

	defer mockServer.Close()

	// Create a test configuration
	config := Config{
		EurekaURL: mockServer.URL,
		App:       "test-app",
		Instance: InstanceData{
			InstanceID: "test-instance",
		},
	}

	// Perform the test
	err := RegisterWithEureka(config, "testuser", "testpassword")
	if err == nil {
		t.Error("Expected an error for unauthorized request, but got none")
	}
}
func TestRegisterWithEurekaInvalidRequestBody(t *testing.T) {
	// Mock Eureka server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send a response with invalid JSON
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid-json"))
	}))

	defer mockServer.Close()

	// Create a test configuration
	config := Config{
		EurekaURL: mockServer.URL,
		App:       "test-app",
		Instance: InstanceData{
			InstanceID: "test-instance",
		},
	}

	// Perform the test
	err := RegisterWithEureka(config, "testuser", "testpassword")
	if err == nil {
		t.Error("Expected an error for invalid JSON response, but got none")
	}
}

func TestRegisterWithEurekaServerError(t *testing.T) {
	// Mock Eureka server for testing
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Send a server error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))

	defer mockServer.Close()

	// Create a test configuration
	config := Config{
		EurekaURL: mockServer.URL,
		App:       "test-app",
		Instance: InstanceData{
			InstanceID: "test-instance",
		},
	}

	// Perform the test
	err := RegisterWithEureka(config, "testuser", "testpassword")
	if err == nil {
		t.Error("Expected an error for server error response, but got none")
	}
}

func TestRegisterWithEurekaRequestError(t *testing.T) {
	// Create a test configuration with an invalid URL
	config := Config{
		EurekaURL: "invalid-url",
		App:       "test-app",
		Instance: InstanceData{
			InstanceID: "test-instance",
		},
	}

	// Perform the test
	err := RegisterWithEureka(config, "testuser", "testpassword")
	if err == nil {
		t.Error("Expected an error for invalid request URL, but got none")
	}
}
