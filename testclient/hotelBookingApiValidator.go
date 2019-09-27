/*
Copyright 2019 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/hotel-booking-api-validator/api"
	"github.com/google/hotel-booking-api-validator/utils"

	pb "github.com/google/hotel-booking-api-validator/v1"
)

var (
	serverAddr           = flag.String("server_addr", "localhost:8080", "Your http server's address in the format of host:port")
	credentialsFile      = flag.String("credentials_file", "", "File containing credentials for your server. Leave blank to bypass authentication. File should have exactly one line of the form 'username:password'.")
	caFile               = flag.String("ca_file", "", "Absolute path to your server's Certificate Authority root cert. Downloading all roots currently recommended by the Google Internet Authority is a suitable alternative https://pki.google.com/roots.pem. Leave blank to connect using http rather than https.")
	fullServerName       = flag.String("full_server_name", "", "Fully qualified domain name. Same name used to sign CN. Only necessary if ca_file is specified and the base URL differs from the server address.")
	availabilityRequest  = flag.String("availability_request", "", "Path to a sample BookingAvailabilityRequest. Format can be either json or pb3")
	submitRequest        = flag.String("submit_request", "", "Path to a sample BookingSubmitRequest. Format can be either json or pb3")
	availabilityEndpoint = flag.String("availability_endpoint", "/v1/BookingAvailability", "URL endpoint for BookingAvailabilityRequest")
	submitEndpoint       = flag.String("submit_endpoint", "/v1/BookingSubmit", "URL endpoint for BookingSubmitRequest")
)

// Stats keep track of the api success and error status
type Stats struct {
	BookingAvailabilitySuccess bool
	BookingSubmitSuccess       bool
}

func logStats(stats Stats) {
	log.Print("\n************* Begin Stats *************\n")
	var totalErrors int

	if *availabilityRequest != "" {
		if stats.BookingAvailabilitySuccess {
			log.Println("BookingAvailability Succeeded")
		} else {
			totalErrors++
			log.Println("BookingAvailability Failed")
		}
	}

	if *submitRequest != "" {
		if stats.BookingSubmitSuccess {
			log.Println("BookingSubmit Succeeded")
		} else {
			totalErrors++
			log.Println("BookingSubmit Failed")
		}
	}

	if stats.BookingSubmitSuccess && stats.BookingAvailabilitySuccess {
		log.Println("All tests pass!")
	}

	log.Print("\n************* End Stats *************\n")
	os.Exit(totalErrors)
}

func main() {
	flag.Parse()
	var stats Stats

	if *availabilityRequest == "" && *submitRequest == "" {
		log.Fatal("You must provide availability_request or submit_request")
	}

	conn, err := api.InitHTTPConnection(*serverAddr, *credentialsFile, *caFile, *fullServerName)
	if err != nil {
		log.Fatalf("Failed to init http connection %v", err)
	}

	if *availabilityRequest != "" {
		utils.LogFlow("Availability Check", "Start")
		// Load search criteria request json/pb from disk
		pbReq := &pb.BookingAvailabilityRequest{}
		if err := utils.LoadRequest(*availabilityRequest, pbReq); err != nil {
			log.Fatalf("Failed to get availability request: %v", err)
		}

		if err = api.BookingAvailability(pbReq, conn, *availabilityEndpoint); err != nil {
			stats.BookingAvailabilitySuccess = false
			log.Printf("Error making BookingAvailabilityRequest: %v", err)
		} else {
			stats.BookingAvailabilitySuccess = true
		}
		utils.LogFlow("Availability Check", "End")
	}

	if *submitRequest != "" {
		utils.LogFlow("Submit Check", "Start")
		// Load search criteria request json/pb from disk
		pbReq := &pb.BookingSubmitRequest{}
		if err := utils.LoadRequest(*submitRequest, pbReq); err != nil {
			log.Fatalf("Failed to get submit request: %v", err)
		}

		if err = api.BookingSubmit(pbReq, conn, *submitEndpoint); err != nil {
			stats.BookingSubmitSuccess = false
			log.Printf("Error making BookingSubmitRequest: %v", err)
		} else {
			stats.BookingSubmitSuccess = true
		}
		utils.LogFlow("Submit Check", "End")
	}
	logStats(stats)
}
