package restlicodec

import "net/url"

type QueryParamsWriter interface {
	Closer
	WriteParams(paramsWriter func(paramNameWriter func(paramName string) Writer) error) error
}

type queryParamsWriter struct {
	*genericWriter
}

func NewQueryParamsWriter() QueryParamsWriter {
	return &queryParamsWriter{genericWriter: newGenericWriter(&urlWriter{stringEscaper: url.QueryEscape})}
}

func (w *queryParamsWriter) WriteParams(paramsWriter func(paramNameWriter func(paramName string) Writer) error) (err error) {
	first := true
	return paramsWriter(func(paramName string) Writer {
		if first {
			first = false
		} else {
			w.RawByte('&')
		}
		w.RawString(paramName)
		w.RawByte('=')
		return w
	})
}
