package restli

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"time"
)

// StandardLogger represents the standard Logger from the log package (log.std for those inclined to read the source
// code)
var StandardLogger stdLogger

type stdLogger struct{}

func (s stdLogger) Printf(format string, v ...interface{}) {
	_ = log.Output(3, fmt.Sprintf(format, v...))
}

type logger interface {
	Printf(format string, v ...interface{})
}

// LoggingRoundTripper is an http.RoundTripper that wraps a backing http.RoundTripper and logs all outbound queries
// (method, URL, headers and body) to the given logger
type LoggingRoundTripper struct {
	http.RoundTripper
	Logger logger
}

func (l *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBytes, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}

	l.Logger.Printf("\n\n%s\n\n\nSubmitting request\n", string(reqBytes))

	startTime := time.Now()
	res, err := l.RoundTripper.RoundTrip(req)
	if err != nil {
		return res, err
	}

	flightTime := time.Since(startTime)

	resBytes, err := httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}

	l.Logger.Printf("Received response in %v:\n\n%s\n\n", flightTime, string(resBytes))

	return res, nil
}
