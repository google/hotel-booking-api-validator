package utils

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var equateErrorMessage cmp.Option = cmp.Comparer(func(x, y error) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	return x.Error() == y.Error()
})

func TestValidateBookingAvailabilityResponse(t *testing.T) {
	data, err := BookingAvailabilityData()
	if err != nil {
		t.Fatalf("error fetching BookingAvailabilityData: %q", err)
	}
	got := ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if got != nil {
		t.Errorf("Expected successful validation, got error %q", got)
	}
}

func TestValidateBookingSubmitResponse(t *testing.T) {
	data, err := BookingSubmitData()
	if err != nil {
		t.Fatalf("error fetching BookingSubmitData: %q", err)
	}
	got := ValidateBookingSubmitResponse(data.ReqPb, data.RespPb)
	if got != nil {
		t.Errorf("Expected successful validation, got error %q", got)
	}
}

func TestValidateBookingSubmitResponseError(t *testing.T) {
	data, err := BookingSubmitData()
	if err != nil {
		t.Fatalf("error fetching BookingSubmitData: %q", err)
	}
	data.RespPb.Reservation.HotelId = "xxx"
	want := fmt.Errorf("echo field(s) did not match request: hotel_id")
	got := ValidateBookingSubmitResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch different value in echo field (diff -got +want): %s", diff)
	}
}

func TestValidateBookingSubmitResponseMissing(t *testing.T) {
	data, err := BookingSubmitData()
	if err != nil {
		t.Fatalf("error fetching BookingSubmitData: %q", err)
	}
	data.RespPb.ApiVersion = 0
	data.RespPb.TransactionId = ""
	data.RespPb.Reservation.Locator.Id = ""
	want := fmt.Errorf("required field(s) missing: api_version, transaction_id, reservation > locator > id")
	got := ValidateBookingSubmitResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch missing required fields (diff -got +want): %s", diff)
	}
}

func TestValidateBookingAvailabilityResponseMissing(t *testing.T) {
	data, err := BookingAvailabilityData()
	if err != nil {
		t.Fatalf("error fetching BookingAvailabilityData: %q", err)
	}
	// default when field not present in json
	data.RespPb.ApiVersion = 0
	data.RespPb.Party.Adults = 0
	data.RespPb.HotelDetails.Address.Address1 = ""
	want := fmt.Errorf("required field(s) missing: api_version, party > adults, hotel_details > address > address1")
	got := ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch missing required fields (diff -got +want): %s", diff)
	}
}

func TestValidateBookingAvailabilityResponseFormat(t *testing.T) {
	data, err := BookingAvailabilityData()
	if err != nil {
		t.Fatalf("error fetching BookingAvailabilityData: %q", err)
	}
	// valid date, but not expected format
	data.ReqPb.StartDate = "20010401"
	data.RespPb.StartDate = "20010401"
	want := fmt.Errorf("error validating format for field(s): start_date")
	got := ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch invalid date format (diff -got +want): %s", diff)
	}
}

func TestValidateBookingAvailabilityResponseArrayValidation(t *testing.T) {
	data, err := BookingAvailabilityData()
	if err != nil {
		t.Fatalf("error fetching BookingAvailabilityData: %q", err)
	}
	// missing room_types > code
	data.RespPb.RoomTypes[1].Code = ""
	want := fmt.Errorf("required field(s) missing: room_types[1] > code")
	got := ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch missing room_type > code (diff -got +want): %s", diff)
	}
}

func TestValidateBookingAvailabilityResponseArrayStructValidation(t *testing.T) {
	data, err := BookingAvailabilityData()
	if err != nil {
		t.Fatalf("error fetching BookingAvailabilityData: %q", err)
	}
	// missing rate_plans > cancellation_policy
	data.RespPb.RatePlans[0].CancellationPolicy = nil
	want := fmt.Errorf("required field(s) missing: rate_plans[0] > cancellation_policy")
	got := ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch missing rate_plans > cancellation_policy (diff -got +want): %s", diff)
	}
}

func TestValidateBookingAvailabilityResponseLineItem(t *testing.T) {
	data, err := BookingAvailabilityData()
	if err != nil {
		t.Fatalf("error fetching BookingAvailabilityData: %q", err)
	}
	// room_rates > line_items > price > amount set to 0
	data.RespPb.RoomRates[0].LineItems[0].Price.Amount = 0
	want := fmt.Errorf("required field(s) missing: room_rates[0] > line_items[0] > price")
	got := ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch price > amount set to 0 (diff -got +want): %s", diff)
	}

	// missing room_rates > line_items > price
	data.RespPb.RoomRates[0].LineItems[0].Price = nil
	want = fmt.Errorf("required field(s) missing: room_rates[0] > line_items[0] > price")
	got = ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch missing room_rates > line_items > price (diff -got +want): %s", diff)
	}
}

func TestValidateBookingAvailabilityResponseInvalidRoomTypeCode(t *testing.T) {
	data, err := BookingAvailabilityData()
	if err != nil {
		t.Fatalf("error fetching BookingAvailabilityData: %q", err)
	}
	// room_rates > room_type_code that does not match any value in room_types > code
	data.RespPb.RoomRates[0].RoomTypeCode = "XXX"
	want := fmt.Errorf("room_rates > room_type_code XXX not present in room_types > code")
	got := ValidateBookingAvailabilityResponse(data.ReqPb, data.RespPb)
	if diff := cmp.Diff(got, want, equateErrorMessage); diff != "" {
		t.Errorf("failed to catch invalid room_type_code (diff -got +want): %s", diff)
	}
}
