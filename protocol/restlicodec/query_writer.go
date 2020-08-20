package restlicodec

import "net/url"

type RestLiQueryParamsWriter interface {
	Closer
	WriteParams(paramsWriter func(paramNameWriter func(paramName string) Writer) error) error
}

type queryParamsWriter struct {
	*genericWriter
}

func NewRestLiQueryParamsWriter() RestLiQueryParamsWriter {
	return &queryParamsWriter{genericWriter: newGenericWriter(&ror2Writer{stringEscaper: url.QueryEscape})}
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
