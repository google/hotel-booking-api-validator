# Hotel Booking - API Validator utility

This repo contains the utility and code samples for validating an implementation
of the Hotel Booking API spec.

The Hotel Booking API Validator utility will send a request body you specify for
any of the service endpoints and validate the response. The validator will
verify that the request and response bodies have valid schema according to the
[v1 API proto file](./proto/v1.proto) and will also validate the response body
correctly echoes the criteria provided in the request body.

Note: This is not an officially supported Google product.

## Test Client

Before using the test utility, the Go programming language must be installed on
your workstation. A precompiled Go binary for your operating system can be
[found here](https://golang.org/dl/).

This guide will assume you're using the default GOPATH and subsequent GOBIN. For
a comprehensive explanation of the GOPATH env variable please see
[this document](https://golang.org/dl/) by the Go team.

### Installing the utilities with Linux

First, build your Go directory structure. A comprehensive guide on the intended
structure of Go code can be [found here](https://golang.org/doc/code.html).

    mkdir -p $HOME/go/bin $HOME/go/pkg $HOME/go/src/github.com/google/hotel-booking-api-validator

Next, add the following to your ~/.bashrc

    export PATH=$PATH:$(go env GOPATH)/bin
    export GOPATH=$(go env GOPATH)
    export GOBIN=$(go env GOPATH)/bin

Source changes

    source ~/.bashrc

Remove any files from a previous installation

    rm -rf $HOME/go/src/github.com/google/hotel-booking-api-validator/

Next, retrieve the utilities from the
[hotel-booking-api-validator repository](https://github.com/google/hotel-booking-api-validator)

    go get github.com/google/hotel-booking-api-validator/...

To install the Hotel Booking API Validator, run

    go install $HOME/go/src/github.com/google/hotel-booking-api-validator/testclient/hotelBookingApiValidator.go

### Installing the utilities with Windows PowerShell

First, build your Go directory structure. A comprehensive guide on the intended
structure of Go code can be [found here](https://golang.org/doc/code.html).

    $env:HOME = $env:USERPROFILE
    md $env:HOME\go\bin
    md $env:HOME\go\pkg
    md $env:HOME\go\src\github.com\google\hotel-booking-api-validator

Next, set the appropriate environment variables

    $env:PATH = $env:PATH + ";" + (go env GOPATH) + "\bin"
    $env:GOPATH = (go env GOPATH)
    $env:GOBIN = (go env GOPATH) + "\bin"

Remove any files from a previous installation

    rd -r $env:HOME\go\src\github.com\google\hotel-booking-api-validator\

Next, retrieve the utilities from the
[hotel-booking-api-validator repository](https://github.com/google/hotel-booking-api-validator)

    cd $env:HOME\go
    # NOTE: You may see output '...cannot find package...' when running the following 'go get' command.
    # These are likely just warnings and you can proceed with the installation.
    go get github.com/google/hotel-booking-api-validator/...

To install the Hotel Booking API Validator, run

    go install $env:HOME\go\src\github.com\google\hotel-booking-api-validator\testclient\hotelBookingApiValidator.go

## Using the validator utility

After following the install steps above an executable should now live under the
path

    $HOME/go/bin/

or

    $env:HOME\go\bin\

All available flags can be displayed using the '--help' flag. The currently
accepted flags are:

```
Usage of hotelBookingApiValidator:
  -server_addr string
        Your http server's address in the format of host:port (default "example.com:80")
  -full_server_name string
        Fully qualified domain name. Same name used to sign CN. Only necessary if ca_file is specified and the base URL differs from the server address.
  -ca_file string
        Absolute path to your server's Certificate Authority root cert. Downloading all roots currently recommended by the Google Internet Authority is a suitable alternative https://pki.google.com/roots.pem. Leave blank to connect using http rather than https.
  -credentials_file string
        File containing credentials for your server. Leave blank to bypass authentication. File should have exactly one line of the form 'username:password'.
  -availability_endpoint string
        URL endpoint for BookingAvailabilityRequest (default "/v1/BookingAvailability")
  -submit_endpoint string
        URL endpoint for BookingSubmitRequest (default "/v1/BookingSubmit")
  -availability_request string
        Path to a sample BookingAvailabilityRequest. Format can be either json or pb3
  -submit_request string
        Path to a sample BookingSubmitRequest. Format can be either json or pb3
```

Example Usage:

```bash
cd $HOME/go
export DATA_PATH=src/github.com/google/hotel-booking-api-validator/data

bin/hotelBookingApiValidator \
  --server_addr=localhost:8080 \
  --availability_request=$DATA_PATH/BookingAvailabilityRequest.json

bin/hotelBookingApiValidator \
  --server_addr=localhost:443 \
  --availability_endpoint=/api/booking/availability
  --availability_request=$DATA_PATH/BookingAvailabilityRequest.json

bin/hotelBookingApiValidator \
  --server_addr=external-dns:443 \
  --ca_file=/path/to/external-dns.pem \
  --availability_request=$DATA_PATH/BookingAvailabilityRequest.json \
  --submit_request=$DATA_PATH/BookingSubmitRequest.json
```

### Sample Request and Response documents

Example json request and response documents for the BookingAvailability service
and BookingSubmit service can be found in the [data](./data/) folder of this
repository. The
[BookingAvailabilityRequest.json](./data/BookingAvailabilityRequest.json) and
[BookingSubmitRequest.json](./data/BookingSubmitRequest.json) files are suitable
as testing input with modifications to match the properties available through
your service.

### Testing

It is important that as part of testing you verify all aspects of the server
invocation - authentication, availability, booking, and error handling.

### Parsing the output

The validation utility will output the logs to stdout. Each line will begin with
a timestamp in RFC3339 format. The output file contains a complete log of all
Requests and Responses sent/received by the testing utility as well as diffs of
the expected response in the event of errors. Similar to a compiler, an overview
of the entire run can be found at the end of the file for user friendly
digestion.
