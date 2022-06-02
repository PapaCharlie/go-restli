package restli

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

func EncodeTunnelledQuery(httpMethod, query string, body []byte) (newBody []byte, headers http.Header) {
	headers = http.Header{}
	headers.Add(MethodOverrideHeader, httpMethod)

	if len(body) > 0 {
		multiPartBody := &bytes.Buffer{}
		// Since bytes.Buffer never returns an error when written to, all errors returned by mpWriter are ignored
		mpWriter := multipart.NewWriter(multiPartBody)

		w, _ := mpWriter.CreatePart(textproto.MIMEHeader{ContentTypeHeader: {FormUrlEncodedContentType}})
		_, _ = w.Write([]byte(query))

		w, _ = mpWriter.CreatePart(textproto.MIMEHeader{ContentTypeHeader: {ApplicationJsonContentType}})
		_, _ = w.Write(body)

		_ = mpWriter.Close()

		newBody = multiPartBody.Bytes()
		headers.Add(ContentTypeHeader, mime.FormatMediaType(MultipartMixedContentType, map[string]string{
			MultipartBoundary: mpWriter.Boundary(),
		}))
	} else {
		newBody = []byte(query)
		headers.Add(ContentTypeHeader, FormUrlEncodedContentType)
	}

	return newBody, headers
}

func DecodeTunnelledQuery(req *http.Request) (err error) {
	getAndDeleteHeader := func(h string) string {
		v := req.Header.Get(h)
		if v == "" {
			return ""
		}
		req.Header.Del(h)
		return v
	}

	tunnelledMethod := getAndDeleteHeader(MethodOverrideHeader)
	if req.Method != http.MethodPost || tunnelledMethod == "" {
		return nil
	}
	req.Method = tunnelledMethod

	if req.URL.RawQuery != "" {
		return fmt.Errorf("go-restli: Invalid request, cannot specify the %q header with a non-empty query",
			MethodOverrideHeader)
	}

	body := new(bytes.Buffer)
	_, err = io.Copy(body, req.Body)
	if err != nil {
		return err
	}
	err = req.Body.Close()
	if err != nil {
		return err
	}
	req.Body = nil

	mediaType, params, err := mime.ParseMediaType(getAndDeleteHeader(ContentTypeHeader))
	switch mediaType {
	case FormUrlEncodedContentType:
		req.URL.RawQuery = body.String()
		req.RequestURI = req.URL.RequestURI()
		req.Body = http.NoBody
		return nil
	case MultipartMixedContentType:
		r := multipart.NewReader(body, params[MultipartBoundary])
		for {
			part, err := r.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			contentType := part.Header.Get(ContentTypeHeader)
			switch contentType {
			case FormUrlEncodedContentType:
				query, _ := io.ReadAll(part)
				req.URL.RawQuery = string(query)
			case ApplicationJsonContentType:
				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, part)
				if err != nil {
					return err
				}
				req.Body = io.NopCloser(buf)
				req.Header.Set(ContentTypeHeader, ApplicationJsonContentType)
			default:
				return fmt.Errorf("go-restli: Unknown tunnelled %s: %q", ContentTypeHeader, contentType)
			}
		}
		if req.URL.RawQuery == "" {
			return fmt.Errorf("go-restli: No query specified in tunneled query request")
		}
		req.RequestURI = req.URL.RequestURI()
		if req.Body == nil {
			return fmt.Errorf("go-restli: No body specified in tunneled query request")
		}
		return nil
	default:
		return nil
	}
}
