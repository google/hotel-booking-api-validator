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
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var reader = ioutil.ReadFile

// LogFlow is a convenience function for logging common flows..
func LogFlow(f string, status string) {
	log.Println(strings.Join([]string{"\n##########\n", status, f, "Flow", "\n##########"}, " "))
}

// LoadRequest loads the request file and returns it's parsed version in pb.
func LoadRequest(fp string, pbReq proto.Message) error {
	content, err := reader(fp)

	if err != nil {
		return fmt.Errorf("unable to read input file: %v", err)
	}
	if path.Ext(fp) == ".json" {
		if err := jsonpb.UnmarshalString(string(content), pbReq); err != nil {
			return fmt.Errorf("unable to parse request as json: %v", err)
		}
		return nil
	}
	if path.Ext(fp) == ".pb3" {
		if err := proto.UnmarshalText(string(content), pbReq); err != nil {
			return fmt.Errorf("unable to parse request as pb3: %v", err)
		}
		return nil
	}
	return fmt.Errorf("unexpected extension for file %q, expected .json or .pb3", fp)
}
