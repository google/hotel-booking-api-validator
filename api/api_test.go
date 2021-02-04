package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmpopts/cmpopts"
	"github.com/golang/protobuf/jsonpb"
	"github.com/google/hotel-booking-api-validator/utils"
)

type ReadFileFunc func(filename string) ([]byte, error)

type FakeFileReader map[string][]byte

func (f FakeFileReader) ReadFile(filename string) ([]byte, error) {
	v := f[filename]
	if v == nil {
		return nil, fmt.Errorf("no record of filename: %s", filename)
	}
	return v, nil
}

func NewFakeFileReader() (FakeFileReader, error) {
	cafile, err := utils.ReadFakeCert()
	if err != nil {
		return FakeFileReader{}, err
	}
	return FakeFileReader{
		"/path/to/credentials": []byte("username:password"),
		"/path/to/pem":         []byte(cafile),
	}, nil
}

func setupMockReader(t *testing.T) {
	r, err := NewFakeFileReader()
	if err != nil {
		t.Fatal(err)
	}
	reader = r.ReadFile
}

func NewFakeHTTPClient(t *testing.T, response string) (*HTTPConnection, *httptest.Server) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, response)
	}))
	return &HTTPConnection{
		client:      server.Client(),
		credentials: "",
		marshaler:   &jsonpb.Marshaler{OrigName: true},
		baseURL:     server.URL,
	}, server
}

func TestBookingAvailability(t *testing.T) {
	data, err := utils.BookingAvailabilityData()
	if err != nil {
		t.Fatal(err)
	}
	conn, server := NewFakeHTTPClient(t, data.Resp)
	defer server.Close()
	if err := BookingAvailability(data.ReqPb, conn, "/BookingAvailability"); err != nil {
		t.Error(err)
	}
}

func TestBookingSubmit(t *testing.T) {
	data, err := utils.BookingSubmitData()
	if err != nil {
		t.Fatal(err)
	}
	conn, server := NewFakeHTTPClient(t, data.Resp)
	defer server.Close()
	if err := BookingSubmit(data.ReqPb, conn, "/BookingSubmit"); err != nil {
		t.Error(err)
	}
}

func TestBookingAvailabilityValidationError(t *testing.T) {
	data, err := utils.BookingAvailabilityData()
	if err != nil {
		t.Fatal(err)
	}
	conn, server := NewFakeHTTPClient(t, data.Resp)
	defer server.Close()
	// Change a value from the request to throw a validation error
	data.ReqPb.HotelId = "xxx"
	want := "Validation error: echo field(s) did not match request: hotel_id"
	if err := BookingAvailability(data.ReqPb, conn, ""); err != nil {
		if err.Error() != want {
			t.Errorf("BookingAvailability(), got [%v] want [%v]", err, want)
		}
	}
}

func TestBookingSubmitValidationError(t *testing.T) {
	data, err := utils.BookingSubmitData()
	if err != nil {
		t.Fatal(err)
	}
	conn, server := NewFakeHTTPClient(t, data.Resp)
	defer server.Close()
	// Change a value from the request to throw a validation error
	data.ReqPb.HotelId = "xxx"
	want := "Validation error: echo field(s) did not match request: hotel_id"
	if err := BookingSubmit(data.ReqPb, conn, ""); err != nil {
		if err.Error() != want {
			t.Errorf("BookingSubmit(), got [%v] want [%v]", err, want)
		}
	}
}

func TestHTTPConnectionURL(t *testing.T) {
	cases := []struct {
		serverAddr      string
		credentialsFile string
		caFile          string
		fullServerName  string
		want            string
	}{
		{
			serverAddr:      "localhost:8080",
			credentialsFile: "",
			caFile:          "",
			fullServerName:  "",
			want:            "http://localhost:8080/test",
		},
		{
			serverAddr:      "localhost:8080",
			credentialsFile: "",
			caFile:          "/path/to/pem",
			fullServerName:  "",
			want:            "https://localhost:8080/test",
		},
	}
	setupMockReader(t)
	for i, tc := range cases {
		conn, err := InitHTTPConnection(tc.serverAddr, tc.credentialsFile, tc.caFile, tc.fullServerName)
		if err != nil {
			t.Errorf("InitHTTPConnection() #%d returned error: %v", i, err)
			continue
		}
		got := conn.getURL("/test")
		if got != tc.want {
			t.Errorf("InitHTTPConnection(), got [%v] want [%v]", got, tc.want)
		}
	}
}

func TestHTTPConnectionCredentials(t *testing.T) {
	cases := []struct {
		serverAddr      string
		credentialsFile string
		caFile          string
		fullServerName  string
		want            string
	}{
		{
			serverAddr:      "localhost:8080",
			credentialsFile: "",
			caFile:          "",
			fullServerName:  "",
			want:            "",
		},
		{
			serverAddr:      "localhost:8080",
			credentialsFile: "/path/to/credentials",
			caFile:          "",
			fullServerName:  "",
			want:            "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
		},
	}
	setupMockReader(t)
	for i, tc := range cases {
		conn, err := InitHTTPConnection(tc.serverAddr, tc.credentialsFile, tc.caFile, tc.fullServerName)
		if err != nil {
			t.Errorf("InitHTTPConnection() #%d returned error: %v", i, err)
			continue
		}
		got := conn.credentials
		if got != tc.want {
			t.Errorf("InitHTTPConnection(), got [%v] want [%v]", got, tc.want)
		}
	}
}

func TestHTTPConnectionCert(t *testing.T) {
	cafile, err := utils.ReadFakeCert()
	if err != nil {
		t.Fatal(err)
	}
	fakeCertPool := x509.NewCertPool()
	fakeCertPool.AppendCertsFromPEM([]byte(cafile))
	cases := []struct {
		serverAddr      string
		credentialsFile string
		caFile          string
		fullServerName  string
		want            *tls.Config
		wantServerName  string
	}{
		{
			serverAddr:      "test1:8080",
			credentialsFile: "",
			caFile:          "",
			fullServerName:  "",
			wantServerName:  "http://test1:8080",
		},
		{
			serverAddr:      "test2:8080",
			credentialsFile: "",
			caFile:          "/path/to/pem",
			fullServerName:  "",
			want: &tls.Config{
				RootCAs:    fakeCertPool,
				ServerName: "",
			},
			wantServerName: "https://test2:8080",
		},
		{
			serverAddr:      "test3.com:8080",
			credentialsFile: "",
			caFile:          "/path/to/pem",
			fullServerName:  "test3.com",
			want: &tls.Config{
				RootCAs:    fakeCertPool,
				ServerName: "test3.com",
			},
			wantServerName: "https://test3.com:8080",
		},
	}
	setupMockReader(t)
	for i, tc := range cases {
		conn, err := InitHTTPConnection(tc.serverAddr, tc.credentialsFile, tc.caFile, tc.fullServerName)
		if err != nil {
			t.Errorf("InitHTTPConnection() #%d returned error: %v", i, err)
			continue
		}

		if conn.baseURL != tc.wantServerName {
			t.Errorf("InitHTTPConnection(%s), got [%v] want [%v]", tc.serverAddr, conn.baseURL, tc.wantServerName)
		}
		if tc.want != nil && !cmp.Equal(conn.config.RootCAs.Subjects(), tc.want.RootCAs.Subjects(), cmpopts.EquateEmpty()) {
			t.Errorf("InitHTTPConnection(%s), got [%v] want [%v]", tc.serverAddr, conn.config.RootCAs, tc.want.RootCAs)
		}
	}
}
