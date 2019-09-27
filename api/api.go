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

// Package api contains validation wrappers over BookingService endpoints.
package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/jsonpb"

	"github.com/google/hotel-booking-api-validator/utils"

	pb "github.com/google/hotel-booking-api-validator/v1"
)

// TimeoutDuration represents the API response timeout duration in miliseconds.
const TimeoutDuration = 30 * time.Second

var reader = ioutil.ReadFile

// HTTPConnection is a convenience struct for holding connection-related objects.
type HTTPConnection struct {
	client      *http.Client
	credentials string
	marshaler   *jsonpb.Marshaler
	baseURL     string
}

// InitHTTPConnection creates and returns a new HTTPConnection object with a given server address and username/password.
func InitHTTPConnection(serverAddr, credentialsFile, caFile, fullServerName string) (*HTTPConnection, error) {
	// Set up username/password.
	credentials, err := setupCredentials(credentialsFile)
	if err != nil {
		return nil, err
	}
	config, err := setupCertConfig(caFile, fullServerName)
	if err != nil {
		return nil, err
	}
	protocol := "http"
	if config != nil {
		protocol = "https"
	}
	return &HTTPConnection{
		client: &http.Client{
			Timeout:   TimeoutDuration,
			Transport: &http.Transport{TLSClientConfig: config},
		},
		credentials: credentials,
		marshaler:   &jsonpb.Marshaler{OrigName: true},
		baseURL:     protocol + "://" + serverAddr,
	}, nil
}

func (h HTTPConnection) getURL(endpoint string) string {
	if endpoint != "" {
		return fmt.Sprintf("%v%v", h.baseURL, endpoint)
	}
	return h.baseURL
}

func setupCredentials(credentialsFile string) (string, error) {
	var credentials string
	if credentialsFile != "" {
		data, err := reader(credentialsFile)
		if err != nil {
			return "", err
		}
		credentials = "Basic " + base64.StdEncoding.EncodeToString([]byte(strings.Replace(string(data), "\n", "", -1)))
	}
	return credentials, nil
}

func setupCertConfig(caFile, fullServerName string) (*tls.Config, error) {
	if caFile == "" {
		return nil, nil
	}
	b, err := reader(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read root certificates file: %v", err)
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, errors.New("failed to parse root certificates, please check your roots file (ca_file flag) and try again")
	}
	return &tls.Config{
		RootCAs:    cp,
		ServerName: fullServerName,
	}, nil
}

func logHTTPRequest(rpcName string, httpReq *http.Request) {
	log.Printf("RPC %s Request. Sent(unix): %s, Url: %s, Method: %s, Header: %s, Body: %v\n", rpcName, time.Now().UTC().Format(time.RFC850), httpReq.URL, httpReq.Method, httpReq.Header, httpReq.Body)
}

func logHTTPResponse(rpcName, bodyString string) {
	log.Printf("RPC %s Response. Received(unix): %s, Response %s\n", rpcName, time.Now().UTC().Format(time.RFC850), bodyString)
}

// sendRequest sets up and sends the relevant HTTP request to the server and returns the HTTP response.
func sendRequest(endpoint, req string, conn *HTTPConnection) (string, error) {
	httpReq, err := http.NewRequest("POST", conn.getURL(endpoint), bytes.NewBuffer([]byte(req)))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", conn.credentials)
	logHTTPRequest(endpoint, httpReq)
	httpResp, err := conn.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("Invalid response. %s yielded error: %v", endpoint, err)
	}
	defer httpResp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return "", fmt.Errorf("Could not read http response body: %v", err)
	}
	bodyString := string(bodyBytes)
	logHTTPResponse(endpoint, bodyString)
	return bodyString, nil
}

// BookingAvailability requests the rooms and metadata, that are available for a specified request context
func BookingAvailability(reqPB *pb.BookingAvailabilityRequest, conn *HTTPConnection, endpoint string) error {
	req, err := conn.marshaler.MarshalToString(reqPB)
	if err != nil {
		return fmt.Errorf("Could not convert pb3 to json: %v, Error: %v", reqPB, err)
	}

	httpResp, err := sendRequest(endpoint, req, conn)
	if err != nil {
		return fmt.Errorf("HTTP response yielded error: %v", err)
	}
	var respPB pb.BookingAvailabilityResponse
	if err := jsonpb.UnmarshalString(httpResp, &respPB); err != nil {
		return fmt.Errorf("Could not parse HTTP response to pb3: %v", err)
	}

	if err := utils.ValidateBookingAvailabilityResponse(reqPB, &respPB); err != nil {
		return fmt.Errorf("Validation error: %v", err)
	}

	return nil
}

// BookingSubmit requests the rooms and metadata, that are available for a specified request context
func BookingSubmit(reqPB *pb.BookingSubmitRequest, conn *HTTPConnection, endpoint string) error {
	req, err := conn.marshaler.MarshalToString(reqPB)
	if err != nil {
		return fmt.Errorf("Could not convert pb3 to json: %v, Error: %v", reqPB, err)
	}

	httpResp, err := sendRequest(endpoint, req, conn)
	if err != nil {
		return fmt.Errorf("%s: HTTP response yielded error: %v", endpoint, err)
	}
	var respPB pb.BookingSubmitResponse
	if err := jsonpb.UnmarshalString(httpResp, &respPB); err != nil {
		return fmt.Errorf("%s: Could not parse HTTP response to pb3: %v", endpoint, err)
	}

	if err := utils.ValidateBookingSubmitResponse(reqPB, &respPB); err != nil {
		return fmt.Errorf("Validation error: %v", err)
	}

	return nil
}
