# Go Eureka Library

[![Build Status](https://travis-ci.org/Anderson-losada-expe/go-eureka-lib.svg?branch=main)](https://travis-ci.org/Anderson-losada-expe/go-eureka-lib)
[![Go Report Card](https://goreportcard.com/badge/github.com/Anderson-losada-expe/go-eureka-lib)](https://goreportcard.com/report/github.com/Anderson-losada-expe/go-eureka-lib)
[![GoDoc](https://godoc.org/github.com/Anderson-losada-expe/go-eureka-lib?status.svg)](https://godoc.org/github.com/Anderson-losada-expe/go-eureka-lib)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Go Eureka Library is a Go package that simplifies the registration of your service instances with Eureka, a service registry.

## Features

- Easy registration of service instances with Eureka.
- Support for custom authentication headers.
- JSON configuration for service instances.

## Installation

```bash
go get -u github.com/Anderson-losada-expe/go-eureka-lib
```

## Usage
- Here's a quick example of how to use the library:

```go
package main

import (
	"fmt"

	"github.com/Anderson-losada-expe/go-eureka-lib"
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
```

## License

This library is licensed under the MIT License - see the LICENSE file for details.
