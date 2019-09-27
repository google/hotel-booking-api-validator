package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"

	pb "github.com/google/hotel-booking-api-validator/v1"
)

type ReadFileFunc func(filename string) ([]byte, error)

type MarshalMessageFunc func(io.Writer, proto.Message) error

type FakeFileReader struct {
	Path     string
	Contents []byte
}

func (f FakeFileReader) ReadFile(filename string) ([]byte, error) {
	return f.Contents, nil
}

type FakeMessageReader struct {
	Path      string
	Message   proto.Message
	Marshaler MarshalMessageFunc
}

func (f FakeMessageReader) ReadFile(filename string) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := f.Marshaler(buf, f.Message); err != nil {
		return nil, err
	}
	return ioutil.ReadAll(buf)
}

func TestLoadRequests(t *testing.T) {
	jsonMarshaler := &jsonpb.Marshaler{}
	testFormats := []struct {
		path      string
		marshaler MarshalMessageFunc
	}{
		{
			path:      "test_request.json",
			marshaler: jsonMarshaler.Marshal,
		},
		{
			path:      "test_request.pb3",
			marshaler: proto.MarshalText,
		},
	}

	availabilityData, err := BookingAvailabilityData()
	if err != nil {
		t.Fatal(err)
	}
	submitData, err := BookingSubmitData()
	if err != nil {
		t.Fatal(err)
	}
	testsMessages := []struct {
		message   proto.Message
		protoType proto.Message
	}{
		{
			message:   availabilityData.ReqPb,
			protoType: &pb.BookingAvailabilityRequest{},
		},
		{
			message:   submitData.ReqPb,
			protoType: &pb.BookingSubmitRequest{},
		},
	}

	for _, tf := range testFormats {
		for _, tm := range testsMessages {
			fake := FakeMessageReader{Path: tf.path, Message: tm.message, Marshaler: tf.marshaler}
			reader = fake.ReadFile
			got := tm.protoType
			err := LoadRequest(tf.path, got)
			if err != nil {
				t.Errorf("LoadRequest(%s, %s) returned an error: %v", tf.path, got, err)
				continue
			}
			if !proto.Equal(got, tm.message) {
				t.Errorf("Failed, got [%v] want [%v]", got, tm.message)
			}
		}
	}
}
