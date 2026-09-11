package config

import (
	"fmt"
	"regexp"
	"strings"
)

func validate() []string {
	var errs []string
	var localerrs []error

	for _, e := range Engine {
		localerrs = localerrs[:0]

		nameValidation(e.Name, &localerrs)
		keyValidation(e.Key, &localerrs)

		for _, err := range localerrs {
			errs = append(errs, fmt.Sprintf("engine %q: %v", e.Name, err))
		}
	}

	for _, v := range Variant {
		localerrs = localerrs[:0]

		nameValidation(v.Name, &localerrs)
		keyValidation(v.Key, &localerrs)
		// TODO: validate v.Match length -1 == v.Args length OR v.Match length == v.Args length.

		for _, err := range localerrs {
			errs = append(errs, fmt.Sprintf("variant %q: %v", v.Name, err))
		}
	}

	localerrs = localerrs[:0]

	if len(Engine) == 0 {
		errs = append(errs, "require at least 1 engine")
	}

	if len(Server.FileTypes) == 0 {
		errs = append(errs, "require at least 1 file type")
	}

	return errs
}

// Regex to check if a string contains only alphanumeric characters (a-z, A-Z, 0-9)
var alphanumericRegex = regexp.MustCompile("^[a-zA-Z0-9]+$")

func keyValidation(value string, errs *[]error) {
	if len(value) == 0 || len(value) > 16 {
		*errs = append(*errs, fmt.Errorf("key length must be between 1 and 16 characters"))
	}

	if value != "" && !alphanumericRegex.MatchString(value) {
		*errs = append(*errs, fmt.Errorf("key must contain only alphanumeric characters"))
	}
}

func nameValidation(value string, errs *[]error) {
	val := strings.TrimSpace(value)
	if val == "" {
		*errs = append(*errs, fmt.Errorf("name cannot be whitespace or blank"))
	}
	if strings.Contains(val, "..") {
		*errs = append(*errs, fmt.Errorf("name cannot contain '..'"))
	}
	if strings.ContainsAny(val, "/\\") {
		*errs = append(*errs, fmt.Errorf("name cannot contain slashes or backslashes"))
	}
}
