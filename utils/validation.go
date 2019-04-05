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
	"fmt"
	"log"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"

	pb "github.com/google/hotel-booking-api-validator/v1"
)

type validationTests []struct {
	field string
	want  interface{}
	got   interface{}
}

func compareFields(validationTests validationTests) error {
	var errorFields []string

	for _, tt := range validationTests {
		if diff := cmp.Diff(tt.got, tt.want, cmp.Comparer(proto.Equal)); diff != "" {
			errorFields = append(errorFields, tt.field)
			log.Println(fmt.Errorf("%s did not match (-got +want)\n%s", tt.field, diff))
		}
	}

	if len(errorFields) > 0 {
		return fmt.Errorf("Error validating field(s): %v", strings.Join(errorFields, ","))
	}

	return nil
}

// ValidateBookingAvailabilityResponse ensures the availability search criteria matches the echoed response.
func ValidateBookingAvailabilityResponse(req *pb.BookingAvailabilityRequest, resp *pb.BookingAvailabilityResponse) error {
	return compareFields(validationTests{
		{"hotel_id", req.GetHotelId(), resp.GetHotelId()},
		{"start_date", req.GetStartDate(), resp.GetStartDate()},
		{"end_date", req.GetEndDate(), resp.GetEndDate()},
		{"party", req.GetParty(), resp.GetParty()},
	})
}

// ValidateBookingSubmitResponse ensures the booking submit criteria matches the echoed response.
func ValidateBookingSubmitResponse(req *pb.BookingSubmitRequest, resp *pb.BookingSubmitResponse) error {
	return compareFields(validationTests{
		{"hotel_id", req.GetHotelId(), resp.GetReservation().GetHotelId()},
		{"start_date", req.GetStartDate(), resp.GetReservation().GetStartDate()},
		{"end_date", req.GetEndDate(), resp.GetReservation().GetEndDate()},
		{"customer", req.GetCustomer(), resp.GetReservation().GetCustomer()},
		{"traveler", req.GetTraveler(), resp.GetReservation().GetTraveler()},
		{"room_rate", req.GetRoomRate(), resp.GetReservation().GetRoomRate()},
	})
}
