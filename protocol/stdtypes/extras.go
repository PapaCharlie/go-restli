package stdtypes

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

func (e *ErrorResponse) Error() string {
	b := strings.Builder{}
	b.WriteString("RestLiError(status: ")

	if e.Status != nil {
		b.WriteString(strconv.Itoa(int(*e.Status)))
	} else {
		b.WriteString("UNKNOWN")
	}

	if e.ExceptionClass != nil {
		b.WriteString(", exceptionClass: ")
		b.WriteString(*e.ExceptionClass)
	}

	if e.Message != nil {
		b.WriteString(", message: ")
		b.WriteString(*e.Message)
	}

	b.WriteString(")")

	return b.String()

}

func (e *ErrorResponse) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		io.WriteString(s, e.Error())
		if e.StackTrace != nil {
			io.WriteString(s, "\n"+*e.StackTrace)
		}
	case 's':
		io.WriteString(s, e.Error())
	}
}

func NewPagingContext(start, count int32) PagingContext {
	return PagingContext{
		Count: &count,
		Start: &start,
	}
}
