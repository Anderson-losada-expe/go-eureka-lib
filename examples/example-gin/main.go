package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Anderson-losada-expe/go-eureka-lib/eureka"
	"github.com/gin-gonic/gin"
)

func main() {
	// Instance configuration
	instanceData := eureka.InstanceData{
		InstanceID: "app-go-gin-instance",
		HostName:   "localhost",
		App:        "app-go-gin",
		IPAddr:     "127.0.0.1",
		Status:     eureka.StatusUp,
		Port: eureka.PortData{
			Value:   8080,
			Enabled: "true",
		},
		HomePageURL:      "http://localhost:8080/",
		StatusPageURL:    "http://localhost:8080/status",
		HealthCheckURL:   "http://localhost:8080/health",
		VIPAddress:       "app-go-gin",
		SecureVIPAddress: "app-go-gin",
		CountryID:        1,
		DataCenterInfo: eureka.DataCenterInfoData{
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
			Name:  "MyOwn",
		},
	}

	// Configuración de Eureka
	eurekaConfig := eureka.Config{
		EurekaURL: "http://127.0.0.1:8761/eureka/apps/app-go-gin",
		App:       "app-go-gin",
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

	// Gin router configuration
	router := gin.Default()

	// Example path to check status
	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	// Example route to check health
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"health": "OK"})
	})

	// Start Gin server
	go func() {
		if err := router.Run(":8080"); err != nil {
			fmt.Println("Error al iniciar el servidor Gin:", err)
		}
	}()

	// Handle closing signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Aplicación terminada. Cerrando...")
}
