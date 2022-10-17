package restlicodec

import (
	"sort"
	"strings"
)

type RestLiQueryParamsWriter interface {
	WriteParams(paramsWriter MapWriter) (string, error)
}

var unescapedQueryCharacters = func() map[byte]struct{} {
	const chars = `!$*-./0123456789;?@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~`
	m := make(map[byte]struct{}, len(chars))
	for i := range chars {
		m[chars[i]] = struct{}{}
	}
	return m
}()

// Ror2QueryEscape query-escapes the given string using Rest.li's query escaper (same as Ror2PathEscape). Using the same
// technique of generating all the UTF8 characters, a handful of characters were found to be escaped differently than
// how a normal query encoder would do it.
func Ror2QueryEscape(s string) string {
	buf := new(strings.Builder)
	for _, c := range []byte(s) {
		if c == ' ' {
			hexEscape(buf, ' ')
		} else if _, ok := unescapedQueryCharacters[c]; ok {
			buf.WriteByte(c)
		} else {
			hexEscape(buf, c)
		}
	}
	return buf.String()
}

func NewRestLiQueryParamsWriter() Writer {
	return newGenericWriter(newRor2Writer(Ror2QueryEscape), nil)
}

func BuildQueryParams(paramsWriter MapWriter) (out string, err error) {
	type entry struct {
		param  string
		writer *genericWriter
	}
	var entries []entry

	err = paramsWriter(func(param string) Writer {
		e := entry{
			param:  param,
			writer: newGenericWriter(newRor2Writer(Ror2QueryEscape), nil),
		}
		entries = append(entries, e)
		return e.writer
	})
	if err != nil {
		return "", err
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].param < entries[j].param })

	builder := new(strings.Builder)
	for i, e := range entries {
		if i != 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(e.param)
		builder.WriteByte('=')

		_, _ = e.writer.getWriter().Buffer.DumpTo(builder)
	}

	return builder.String(), nil
}
