package tests

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	reqBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Could not read %s", filename)
	}
	r := bufio.NewReader(bytes.NewBuffer(reqBytes))

	req, err := http.ReadRequest(r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not read request")
	}
	// ReadRequest only reads the leading HTTP protocol bytes (e.g. GET /foo HTTP/1.1) and the headers. What remains of
	// the buffer is the body of the request, which we need to preserve for subsequent reads
	reqBytes, _ = ioutil.ReadAll(r)
	return req, adjustContentLength(filename, reqBytes, req.Header), nil
}

func ReadResponseFromFile(filename string, req *http.Request) (*http.Response, []byte, error) {
	resBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "Could not read %s", filename)
	}
	r := bufio.NewReader(bytes.NewBuffer(resBytes))

	res, err := http.ReadResponse(r, req)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not read request")
	}
	// ReadResponse only reads the leading HTTP protocol bytes (e.g. GET /foo HTTP/1.1) and the headers. What remains of
	// the buffer is the body of the response, which we need to preserve for subsequent reads
	resBytes, _ = ioutil.ReadAll(r)
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
