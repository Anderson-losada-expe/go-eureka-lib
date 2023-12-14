package main

import (
	"Anderson-losada-expe/go-eureka-lib/eureka"
	"fmt"
)

func main() {
	// Instance configuration
	instanceData := eureka.InstanceData{
		InstanceID: "app-go-examples-instance",
		HostName:   "localhost",
		App:        "app-go-examples",
		IPAddr:     "127.0.0.1",
		Status:     eureka.StatusUp,
		Port: eureka.PortData{
			Value:   10000,
			Enabled: "true",
		},
		SecurePort: eureka.PortData{
			Value:   443,
			Enabled: "false",
		},
		HomePageURL:      "http://localhost:10000/",
		StatusPageURL:    "http://localhost:10000/actuator/info",
		HealthCheckURL:   "http://localhost:10000/actuator/health",
		VIPAddress:       "app-go-examples",
		SecureVIPAddress: "app-go-examples",
		CountryID:        1,
		DataCenterInfo: eureka.DataCenterInfoData{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  "MyOwn", // Production
		},
	}

	// Eureka Configuration
	eurekaConfig := eureka.Config{
		EurekaURL: "http://127.0.0.1:8761/eureka/apps/app-go-example",
		App:       "app-go-examples",
		Instance:  instanceData,
	}

	// Register in Eureka
	err := eureka.RegisterWithEureka(eurekaConfig, "admin", "admin")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Keep the program running to simulate the registered instance
	fmt.Println("Press Ctrl+C to exit.")
	select {}
}
