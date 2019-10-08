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
	"reflect"
	"regexp"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"

	pb "github.com/google/hotel-booking-api-validator/v1"
)

// ISO3166 provides the regular expression for validating a two-letter country code defined by ISO 3166-1
const ISO3166 = `^A[^ABCHJKNPVY]|B[^CKPUX]|C[^BEJPQST]|D[EJKMOZ]|E[CEGHRST]|F[IJKM
	OR]|G[^CJKOVXZ]|H[KMNRTU]|I[DEL-OQ-T]|J[EMOP]|K[EGHIMNPRWYZ]|L[ABCIKR-VY]|M[^B
	IJ]|N[ACEFGILOPRUZ]|OM|P[AE-HKNRSTWY]|QA|R[EOSUW]|S[^FPQUW]|T[^ABEIPQSUXY]|U[A
	GMSYZ]|V[ACEGINU]|WF|WS|YE|YT|Z[AMW]$`

// DateFormat provides the regular expression for validating a date in YYYY-MM-DD format
const DateFormat = `^([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))$`

type validationTest struct {
	field string
	want  interface{}
	got   interface{}
}

type requiredTest struct {
	field string
	got   interface{}
}

type formatTest struct {
	field   string
	value   string
	pattern string
}

// compareFields will ensure each validationTest got and want proto values are equal
func compareFields(v []validationTest) error {
	var errorFields []string

	for _, vv := range v {
		if diff := cmp.Diff(vv.got, vv.want, cmp.Comparer(proto.Equal)); diff != "" {
			errorFields = append(errorFields, vv.field)
			log.Println(fmt.Errorf("%s did not match (-got +want)\n%s", vv.field, diff))
		}
	}

	if len(errorFields) > 0 {
		return fmt.Errorf("echo field(s) did not match request: %v", strings.Join(errorFields, ","))
	}

	return nil
}

// checkRequired will ensure each requiredTest value is not equal to the unsetValue
func checkRequired(r []requiredTest) error {
	var errorFields []string

	for _, rr := range r {
		if reflect.ValueOf(rr.got).IsZero() {
			errorFields = append(errorFields, rr.field)
			log.Println(fmt.Errorf("Required field %s was not set", rr.field))
		}
	}

	if len(errorFields) > 0 {
		return fmt.Errorf("required field(s) missing: %v", strings.Join(errorFields, ", "))
	}

	return nil
}

// validateFormat will ensure each formatTest value matches given pattern
func validateFormat(f []formatTest) error {
	var errorFields []string

	for _, ff := range f {
		matched, err := regexp.Match(ff.pattern, []byte(ff.value))
		if err != nil {
			return err
		}
		if !matched {
			errorFields = append(errorFields, ff.field)
			log.Println(fmt.Errorf("Field %s value %s did not match pattern %v", ff.field, ff.value, ff.pattern))
		}
	}

	if len(errorFields) > 0 {
		return fmt.Errorf("error validating format for field(s): %s", strings.Join(errorFields, ", "))
	}

	return nil
}

// valuePresent will check if value v is present in slice s
func valuePresent(v string, s []string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}

// ValidateBookingAvailabilityResponse ensures the availability search criteria matches the echoed response.
func ValidateBookingAvailabilityResponse(req *pb.BookingAvailabilityRequest, resp *pb.BookingAvailabilityResponse) error {
	// Validate the required fields are present and not set to the default value
	if err := checkRequired([]requiredTest{
		{"api_version", resp.GetApiVersion()},
		{"transaction_id", resp.GetTransactionId()},
		{"hotel_id", resp.GetHotelId()},
		{"party > adults", resp.GetParty().GetAdults()},
		{"hotel_details > name", resp.GetHotelDetails().GetName()},
		{"hotel_details > address > address1", resp.GetHotelDetails().GetAddress().GetAddress1()},
		{"hotel_details > address > city", resp.GetHotelDetails().GetAddress().GetCity()},
		{"hotel_details > address > province", resp.GetHotelDetails().GetAddress().GetProvince()},
	}); err != nil {
		return err
	}
	// Ensure certain fields match expected format
	if err := validateFormat([]formatTest{
		{"start_date", resp.GetStartDate(), DateFormat},
		{"end_date", resp.GetEndDate(), DateFormat},
		{"hotel_details > address > country", resp.GetHotelDetails().GetAddress().GetCountry(), ISO3166},
	}); err != nil {
		return err
	}
	// Ensure response echo fields match request values
	if err := compareFields([]validationTest{
		{"hotel_id", req.GetHotelId(), resp.GetHotelId()},
		{"start_date", req.GetStartDate(), resp.GetStartDate()},
		{"end_date", req.GetEndDate(), resp.GetEndDate()},
		{"party", req.GetParty(), resp.GetParty()},
	}); err != nil {
		return err
	}

	roomTypeCodes := make([]string, len(resp.GetRoomTypes()))
	ratePlanCodes := make([]string, len(resp.GetRatePlans()))

	// Validate each Room Type
	for i, r := range resp.GetRoomTypes() {
		roomTypeCodes[i] = r.GetCode()
		err := checkRequired([]requiredTest{
			{fmt.Sprintf("room_types[%d] > code", i), r.GetCode()},
			{fmt.Sprintf("room_types[%d] > name", i), r.GetName().String()},
		})
		if err != nil {
			return err
		}
	}

	// Validate each Rate Plan
	for i, r := range resp.GetRatePlans() {
		ratePlanCodes[i] = r.GetCode()
		err := checkRequired([]requiredTest{
			{fmt.Sprintf("rate_plans[%d] > code", i), r.GetCode()},
			{fmt.Sprintf("rate_plans[%d] > name", i), r.GetName().String()},
			{fmt.Sprintf("rate_plans[%d] > cancellation_policy", i), r.GetCancellationPolicy()},
		})
		if err != nil {
			return err
		}
	}

	// Validate each Room Rate & ensure room_type_codes and rate_plan_codes exist in response
	for i, r := range resp.GetRoomRates() {
		rt := make([]requiredTest, len(r.GetLineItems()))
		for j, l := range r.GetLineItems() {
			// Ensure price is not zero or unset
			rt[j] = requiredTest{fmt.Sprintf("room_rates[%d] > line_items[%d] > price", i, j), l.GetPrice().GetAmount()}
		}
		rt = append(rt, requiredTest{fmt.Sprintf("room_rates[%d] > code", i), r.GetCode()})
		if err := checkRequired(rt); err != nil {
			return err
		}
		if !valuePresent(r.GetRoomTypeCode(), roomTypeCodes) {
			return fmt.Errorf("room_rates > room_type_code %v not present in room_types > code", r.GetRoomTypeCode())
		}
		if !valuePresent(r.GetRatePlanCode(), ratePlanCodes) {
			return fmt.Errorf("room_rates > rate_plan_code %v not present in rate_plans > code", r.GetRatePlanCode())
		}
	}

	return nil
}

// ValidateBookingSubmitResponse checks for required fields, formats, and matching echo responses.
func ValidateBookingSubmitResponse(req *pb.BookingSubmitRequest, resp *pb.BookingSubmitResponse) error {
	// Validate required fields are present and not set to the default value
	if err := checkRequired([]requiredTest{
		{"api_version", resp.GetApiVersion()},
		{"transaction_id", resp.GetTransactionId()},
		{"status", resp.GetStatus().String()},
		{"reservation > locator > id", resp.GetReservation().GetLocator().GetId()},
	}); err != nil {
		return err
	}

	// Ensure echo response fields match request values
	if err := compareFields([]validationTest{
		{"hotel_id", req.GetHotelId(), resp.GetReservation().GetHotelId()},
		{"start_date", req.GetStartDate(), resp.GetReservation().GetStartDate()},
		{"end_date", req.GetEndDate(), resp.GetReservation().GetEndDate()},
		{"customer", req.GetCustomer(), resp.GetReservation().GetCustomer()},
		{"traveler", req.GetTraveler(), resp.GetReservation().GetTraveler()},
		{"room_rate", req.GetRoomRate(), resp.GetReservation().GetRoomRate()},
	}); err != nil {
		return err
	}

	return nil
}
