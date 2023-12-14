package main

import (
	"Anderson-losada-expe/go-eureka-lib/eureka"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
)

func main() {
	// Instance configuration for Eureka
	instanceData := eureka.InstanceData{
		InstanceID:       "examples-instance",
		HostName:         "localhost",
		App:              "examples-app",
		IPAddr:           "127.0.0.1",
		Status:           eureka.StatusUp,
		Port:             eureka.PortData{Value: 8080, Enabled: "true"},
		SecurePort:       eureka.PortData{Value: 8443, Enabled: "false"},
		HomePageURL:      "http://localhost:8080/",
		StatusPageURL:    "http://localhost:8080/actuator/info",
		HealthCheckURL:   "http://localhost:8080/actuator/health",
		VIPAddress:       "examples-app",
		SecureVIPAddress: "examples-app",
		CountryID:        1,
		DataCenterInfo: eureka.DataCenterInfoData{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  "MyOwn",
		},
	}

	// Eureka Configuration
	eurekaConfig := eureka.Config{
		EurekaURL: "http://127.0.0.1:8761/eureka/apps/example-app",
		App:       "examples-app",
		Instance:  instanceData,
	}

	// Start Echo server
	e := echo.New()

	// Test route to verify server life
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is up and running!")
	})

	// Start logging in Eureka in a goroutine
	go func() {
		err := eureka.RegisterWithEureka(eurekaConfig, "admin", "admin")
		if err != nil {
			fmt.Println("Error registrando en Eureka:", err)
		}
	}()

	// Handle signals to shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("Shutting down the server...")

		// Deregister in Eureka before closing
		err := eureka.DeregisterWithEureka(eurekaConfig, "admin", "admin")
		if err != nil {
			return
		}

		os.Exit(0)
	}()

	// Start the server
	if err := e.Start(":8080"); err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}
