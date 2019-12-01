package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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
	r := bufio.NewReader(f)

	req, err := http.ReadRequest(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not read request")
	}
	// ReadRequest only reads the leading HTTP protocol bytes (e.g. GET /foo HTTP/1.1) and the headers. What remains of
	// the buffer is the body of the request, which we need to preserve for subsequent reads
	reqBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to read full request")
	}
	return req, adjustContentLength(filename, reqBytes, req.Header), nil
}

func ReadResponseFromFile(filename string, req *http.Request) (*http.Response, []byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Panicln("Could not open", filename, err)
	}
	defer f.Close()
	r := bufio.NewReader(f)

	res, err := http.ReadResponse(r, req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not read request")
	}
	// ReadResponse only reads the leading HTTP protocol bytes (e.g. GET /foo HTTP/1.1) and the headers. What remains of
	// the buffer is the body of the response, which we need to preserve for subsequent reads
	resBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to read full request")
	}
	return res, adjustContentLength(filename, resBytes, res.Header), nil
}

func adjustContentLength(filename string, b []byte, h http.Header) []byte {
	const contentLength = "Content-Length"
	b = bytes.Trim(b, "\r\n")
	cl := h.Get(contentLength)
	if cl != "" {
		cli, _ := strconv.Atoi(cl)
		if len(b) != cli {
			log.Printf("Content-Length header in %s indicates %d bytes, but body was %d bytes", filename, cli, len(b))
			h.Set(contentLength, fmt.Sprintf("%d", len(b)))
		}
	}
	return b
}
