package gferr

import (
	"fmt"
	"strings"
)

type MultiError []error

func (e MultiError) Error() string {
	var sb strings.Builder
	sb.WriteRune('[')
	for idx, err := range e {
		if idx > 0 {
			sb.WriteString(", ")
		}
		if err == nil {
			sb.WriteString("<nil>")
			continue
		}
		sb.WriteString(fmt.Sprintf("(%v)", err))
	}
	sb.WriteRune(']')
	return sb.String()
}

func (e MultiError) Unwrap() error {
	if len(e) == 0 {
		return nil
	}
	return e
}
