package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

func Must(req *http.Request, err error) *http.Request {
	if err != nil {
		log.Panicln(err)
	}
	return req
}

func ReadRequestFromFile(filename string) (*http.Request, []byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Panicln("Could not open", filename, err)
	}
	defer f.Close()
	return ReadRequest(bufio.NewReader(f))
}

func ReadRequest(reader *bufio.Reader) (*http.Request, []byte, error) {
	req, err := http.ReadRequest(reader)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not read request")
	}
	// ReadRequest only reads the leading HTTP protocol bytes (e.g. GET /foo HTTP/1.1) and the headers. What remains of
	// the buffer is the body of the request, which we need to preserve for subsequent reads
	reqBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to read full request")
	}
	return req, reqBytes, nil
}

func ReadResponseFromFile(filename string, req *http.Request) (*http.Response, []byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Panicln("Could not open", filename, err)
	}
	defer f.Close()
	return ReadResponse(bufio.NewReader(f), req)
}

func ReadResponse(reader *bufio.Reader, req *http.Request) (*http.Response, []byte, error) {
	res, err := http.ReadResponse(reader, req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not read request")
	}
	// ReadResponse only reads the leading HTTP protocol bytes (e.g. GET /foo HTTP/1.1) and the headers. What remains of
	// the buffer is the body of the response, which we need to preserve for subsequent reads
	resBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to read full request")
	}
	return res, resBytes, nil
}
