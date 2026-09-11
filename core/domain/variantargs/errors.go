package variantargs

import (
	"fmt"
	"strings"
)

type ParseError struct {
	Operation string
	Field     string
	Value     string
	Reason    string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("variantargs: operation %q, field %q, value %q: %s", e.Operation, e.Field, e.Value, e.Reason)
}

type ParseErrors struct {
	Errors []error
}

func (e *ParseErrors) Error() string {
	messages := make([]string, 0, len(e.Errors))
	for _, err := range e.Errors {
		messages = append(messages, err.Error())
	}
	return strings.Join(messages, "; ")
}

func (e *ParseErrors) Unwrap() []error {
	return e.Errors
}
