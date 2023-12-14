package eureka

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Status represents the status of the service instance.
type Status string

const (
	StatusUp           Status = "UP"
	StatusDown         Status = "DOWN"
	StatusOutOfService Status = "OUT_OF_SERVICE"
)

// Config contains the configuration for registering with Eureka.
type Config struct {
	EurekaURL  string
	App        string
	Instance   InstanceData
	AuthHeader string // Custom authorization header
}

// InstanceData contains the data for the instance.
type InstanceData struct {
	InstanceID       string             `json:"instanceId"`
	HostName         string             `json:"hostName"`
	App              string             `json:"app"`
	IPAddr           string             `json:"ipAddr"`
	Status           Status             `json:"status"`
	Port             PortData           `json:"port"`
	SecurePort       PortData           `json:"securePort"`
	HomePageURL      string             `json:"homePageUrl"`
	StatusPageURL    string             `json:"statusPageUrl"`
	HealthCheckURL   string             `json:"healthCheckUrl"`
	VIPAddress       string             `json:"vipAddress"`
	SecureVIPAddress string             `json:"secureVipAddress"`
	CountryID        int                `json:"countryId"`
	DataCenterInfo   DataCenterInfoData `json:"dataCenterInfo"`
}

// PortData contains the port data.
type PortData struct {
	Value   int    `json:"$"`
	Enabled string `json:"@enabled"`
}

// DataCenterInfoData contains the data center data.
type DataCenterInfoData struct {
	Class string `json:"@class"`
	Name  string `json:"name"`
}

// RegisterWithEureka registers the instance with Eureka.
func RegisterWithEureka(config Config, username, password string) error {
	// Convertir a JSON
	jsonData, err := convertToJSON(config.Instance)
	if err != nil {
		return fmt.Errorf("error converting to JSON: %v", err)
	}

	// Encode credentials in "user:password" format in base64
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// Create a custom HTTP request with authorization header
	req, err := createHTTPRequest(config.EurekaURL, jsonData, auth, config.AuthHeader)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Make the registration HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error when making HTTP request: %v", err)
	}
	defer closeResponseBody(resp.Body)

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error registering in Eureka. Status code: %d", resp.StatusCode)
	}

	fmt.Println("Registered in Eureka successfully. 😀👏")
	return nil
}

func convertToJSON(data interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{"instance": data})
}

func createHTTPRequest(url string, body []byte, authHeader, customAuthHeader string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// If a custom authorization header is provided, use it
	if customAuthHeader != "" {
		req.Header.Set("Authorization", customAuthHeader)
	} else {
		req.Header.Set("Authorization", "Basic "+authHeader)
	}

	return req, nil
}

func closeResponseBody(body io.ReadCloser) {
	if body != nil {
		_ = body.Close()
	}
}

// DeregisterWithEureka deregisters the Eureka instance.
func DeregisterWithEureka(config Config, username, password string) error {
	// Encode credentials in "user:password" format in base64
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

	// Create a custom HTTP request with authorization header
	req, err := http.NewRequest("DELETE", config.EurekaURL+"/"+config.Instance.InstanceID, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Basic "+auth)

	// Make the HTTP deregistration request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer closeResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deregistering in Eureka. Status code: %d", resp.StatusCode)
	}

	fmt.Println("Deregistered in Eureka successfully")
	return nil
}
