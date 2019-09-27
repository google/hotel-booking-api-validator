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

package utils

import (
	"io/ioutil"
	"path/filepath"
	"runtime"

	
	jsonpb "github.com/golang/protobuf/jsonpb"
	pb "github.com/google/hotel-booking-api-validator/v1"
)




// basepath is the root directory of this package.
var basepath string

func init() {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath = filepath.Dir(currentFile)
}

func absPath(rel string) string {
	if filepath.IsAbs(rel) {
		return rel
	}

	return filepath.Join(basepath, "../data/", rel)
}

// BookingAvailabilityDataStruct struct
type BookingAvailabilityDataStruct struct {
	Req    string
	ReqPb  *pb.BookingAvailabilityRequest
	Resp   string
	RespPb *pb.BookingAvailabilityResponse
}

// BookingSubmitDataStruct struct
type BookingSubmitDataStruct struct {
	Req    string
	ReqPb  *pb.BookingSubmitRequest
	Resp   string
	RespPb *pb.BookingSubmitResponse
}

func readTestDataFile(filename string) (string, error) {
	f, err := ioutil.ReadFile(absPath(filename))
	if err != nil {
		return "", err
	}
	return string(f), nil
}

// BookingAvailabilityData provides both json string and pb request / response for BookingAvailability
func BookingAvailabilityData() (*BookingAvailabilityDataStruct, error) {
	req, err := readTestDataFile("BookingAvailabilityRequest.json")
	var reqPb pb.BookingAvailabilityRequest
	if err := jsonpb.UnmarshalString(req, &reqPb); err != nil {
		return nil, err
	}
	resp, err := readTestDataFile("BookingAvailabilityResponse.json")
	var respPb pb.BookingAvailabilityResponse
	if err := jsonpb.UnmarshalString(resp, &respPb); err != nil {
		return nil, err
	}
	return &BookingAvailabilityDataStruct{req, &reqPb, resp, &respPb}, err
}

// BookingSubmitData provides both json string and pb request / response for BookingSubmit
func BookingSubmitData() (*BookingSubmitDataStruct, error) {
	req, err := readTestDataFile("BookingSubmitRequest.json")
	var reqPb pb.BookingSubmitRequest
	if err := jsonpb.UnmarshalString(req, &reqPb); err != nil {
		return nil, err
	}
	resp, err := readTestDataFile("BookingSubmitResponse.json")
	var respPb pb.BookingSubmitResponse
	if err := jsonpb.UnmarshalString(resp, &respPb); err != nil {
		return nil, err
	}
	return &BookingSubmitDataStruct{req, &reqPb, resp, &respPb}, err
}

// ReadFakeCert will read the testing cafile.pem from dataDir
func ReadFakeCert() (string, error) {
	f, err := ioutil.ReadFile(absPath("../data/cafile.pem"))
	if err != nil {
		return "", err
	}
	return string(f), nil
}
