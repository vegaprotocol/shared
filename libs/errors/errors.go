package errors

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrIsRequired             = errors.New("is required")
	ErrMustBeValidDate        = errors.New("must be a RFC3339 date")
	ErrMustBePositive         = errors.New("must be positive")
	ErrMustBePositiveOrZero   = errors.New("must be positive or zero")
	ErrMustBeNegative         = errors.New("must be negative")
	ErrMustBeNegativeOrZero   = errors.New("must be negative or zero")
	ErrIsNotValid             = errors.New("is not a valid value")
	ErrIsNotValidNumber       = errors.New("is not a valid number")
	ErrIsNotSupported         = errors.New("is not supported")
	ErrIsUnauthorised         = errors.New("is unauthorised")
	ErrDoesNotMatch           = errors.New("does not match")
	ErrNotAValidInteger       = errors.New("not a valid integer")
	ErrNotAValidFloat         = errors.New("not a valid float")
	ErrServerResponseNone     = errors.New("no response from server")
	ErrServerResponseEmpty    = errors.New("empty response from server")
	ErrServerResponseReadFail = errors.New("failed to read response from server")
	ErrRequestBodyReadFail    = errors.New("failed to read request body")
	ErrInterrupted            = errors.New("interrupted")
	ErrNil                    = errors.New("nil pointer")
	ErrUnrecognisedAction     = errors.New("unrecognised action")
	ErrUnsupportedScheme      = errors.New("unsupported scheme")
	// ErrConnectionNotReady indicated that the network connection to the gRPC server is not ready.
	ErrConnectionNotReady = errors.New("gRPC connection not ready")
	// ErrMissingEmptyConfigSection indicates that a required config file section is missing (not present) or empty (zero-length).
	ErrMissingEmptyConfigSection = errors.New("config file section is missing/empty")
)

func MutuallyExclusiveError(n1, n2 string) error {
	return fmt.Errorf("%s and %s are mutually exclusive", n1, n2)
}

func MustBeSpecifiedError(name string) error {
	return fmt.Errorf("%s must be specified", name)
}

func RequireLessThanError(n, oth string) error {
	return fmt.Errorf("%s must be less than %s", n, oth)
}

func RequireLessThanOrEqualError(n, oth string) error {
	return fmt.Errorf("%s must be less or equal than %s", n, oth)
}

func RequireGreaterThanError(n, oth string) error {
	return fmt.Errorf("%s must be greater than %s", n, oth)
}

func RequireGreaterThanOrEqualError(n, oth string) error {
	return fmt.Errorf("%s must be greater than %s", n, oth)
}

func RequireBetweenValuesError(n, leftInclusive, rightExclusive string) error {
	return fmt.Errorf("%s must be located between %s (inclusive) and %s (exclusive)", n, leftInclusive, rightExclusive)
}

func MustSpecifiedOneOfError(values ...string) error {
	if len(values) < 2 {
		panic("provide at least 2 values")
	}

	lastIndex := len(values) - 1
	firstSegment := strings.Join(values[0:lastIndex], ", ")
	fullSegment := strings.Join([]string{firstSegment, values[lastIndex]}, ", or ")
	return fmt.Errorf("must specified one of %s", fullSegment)
}

func InvalidFormatError(name string) error {
	return fmt.Errorf("%s has not a valid format", name)
}

func UnsupportedValueError(name string, unsupported interface{}, supported []interface{}) error {
	if len(supported) < 2 {
		panic("provide at least 2 supported values")
	}

	supportedFmt := make([]string, 0, len(supported))
	for _, s := range supported {
		supportedFmt = append(supportedFmt, fmt.Sprintf("%v", s))
	}

	lastIndex := len(supportedFmt) - 1
	firstSegment := strings.Join(supportedFmt[0:lastIndex], ", ")
	fullSegment := strings.Join([]string{firstSegment, supportedFmt[lastIndex]}, ", and ")

	return fmt.Errorf("%s does not support value %v, only %s", name, unsupported, fullSegment)
}

func MustBase64EncodedError(name string) error {
	return fmt.Errorf("%s must be base64-encoded", name)
}
