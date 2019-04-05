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

// Package utils contains helper methods for the api validation utility.
package utils

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	pb "github.com/google/hotel-booking-api-validator/v1"
)

// setupCredentials will create a valid http basic auth header value if a valid credentials file is provided.
func setupCredentials(credentialsFile string) (string, error) {
	var credentials string
	if credentialsFile != "" {
		data, err := ioutil.ReadFile(credentialsFile)
		if err != nil {
			return "", err
		}
		credentials = "Basic " + base64.StdEncoding.EncodeToString([]byte(strings.Replace(string(data), "\n", "", -1)))
	}
	return credentials, nil
}

// setupCertConfig attempts to construct a tls config if a valid ca file is provided.
func setupCertConfig(caFile, fullServerName string) (*tls.Config, error) {
	if caFile == "" {
		return nil, nil
	}
	b, err := ioutil.ReadFile(caFile)
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

// LogFlow is a convenience function for logging common flows..
func LogFlow(f string, status string) {
	log.Println(strings.Join([]string{"\n##########\n", status, f, "Flow", "\n##########"}, " "))
}

// LoadAvailabilityRequest loads the request file and returns it's parsed version in pb.
func LoadAvailabilityRequest(fp string) (*pb.BookingAvailabilityRequest, error) {
	var pbReq pb.BookingAvailabilityRequest
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %v", err)
	}
	if path.Ext(fp) == ".json" {
		if err := jsonpb.UnmarshalString(string(content), &pbReq); err != nil {
			return nil, fmt.Errorf("unable to parse request as json: %v", err)
		}
		return &pbReq, nil
	}
	if path.Ext(fp) == ".pb3" {
		if err := proto.Unmarshal(content, &pbReq); err != nil {
			return nil, fmt.Errorf("unable to parse request as pb3: %v", err)
		}
		return &pbReq, nil
	}
	return nil, fmt.Errorf("unexpected extension for file %q, expected .json or .pb3", fp)
}

// LoadSubmitRequest loads the request file and returns it's parsed version in pb.
func LoadSubmitRequest(fp string) (*pb.BookingSubmitRequest, error) {
	var pbReq pb.BookingSubmitRequest
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("unable to read input file: %v", err)
	}
	if path.Ext(fp) == ".json" {
		if err := jsonpb.UnmarshalString(string(content), &pbReq); err != nil {
			return nil, fmt.Errorf("unable to parse request as json: %v", err)
		}
		return &pbReq, nil
	}
	if path.Ext(fp) == ".pb3" {
		if err := proto.Unmarshal(content, &pbReq); err != nil {
			return nil, fmt.Errorf("unable to parse request as pb3: %v", err)
		}
		return &pbReq, nil
	}
	return nil, fmt.Errorf("unexpected extension for file %q, expected .json or .pb3", fp)
}
