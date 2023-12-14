package main

import (
	"Anderson-losada-expe/go-eureka-lib/eureka"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Instance configuration
	instanceData := eureka.InstanceData{
		InstanceID: "app-go-native-instance",
		HostName:   "localhost",
		App:        "app-go-native",
		IPAddr:     "127.0.0.1",
		Status:     eureka.StatusUp,
		Port: eureka.PortData{
			Value:   8080,
			Enabled: "true",
		},
		HomePageURL:      "http://localhost:8080/",
		StatusPageURL:    "http://localhost:8080/status",
		HealthCheckURL:   "http://localhost:8080/health",
		VIPAddress:       "app-go-native",
		SecureVIPAddress: "app-go-native",
		CountryID:        1,
		DataCenterInfo: eureka.DataCenterInfoData{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  "MyOwn",
		},
	}

	// Eureka Configuration
	eurekaConfig := eureka.Config{
		EurekaURL: "http://127.0.0.1:8761/eureka/apps/app-go-native",
		App:       "app-go-native",
		Instance:  instanceData,
	}

	// Register in Eureka when starting the application
	err := eureka.RegisterWithEureka(eurekaConfig, "admin", "admin")
	if err != nil {
		fmt.Println("Error al registrar en Eureka:", err)
		return
	}
	defer func() {
		// Deregister from Eureka when closing the application
		err := eureka.DeregisterWithEureka(eurekaConfig, "admin", "admin")
		if err != nil {
			fmt.Println("Error al derregistrar de Eureka:", err)
		}
	}()

	// Configure route handlers
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "OK"}`))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"health": "OK"}`))
	})

	// Start HTTP server
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Error al iniciar el servidor HTTP:", err)
		}
	}()

	// Handle closing signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Aplicación terminada. Cerrando...")
}
