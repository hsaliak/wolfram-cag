package errs

import (
	"errors"
	"fmt"
	"net"
)

var ErrMissingAPIKey = errors.New("missing API key: pass --api-key or set WOLFRAM_APP_ID")

type InvalidArgsError struct {
	Msg string
}

func (e InvalidArgsError) Error() string {
	return e.Msg
}

type HTTPStatusError struct {
	Code int
	Body string
}

func (e HTTPStatusError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("http status %d", e.Code)
	}
	return fmt.Sprintf("http status %d: %s", e.Code, e.Body)
}

type DecodeError struct {
	Err error
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("decode error: %v", e.Err)
}

func (e DecodeError) Unwrap() error {
	return e.Err
}

type NetworkError struct {
	Err error
}

func (e NetworkError) Error() string {
	return fmt.Sprintf("network error: %v", e.Err)
}

func (e NetworkError) Unwrap() error {
	return e.Err
}

type TimeoutError struct {
	Err error
}

func (e TimeoutError) Error() string {
	return fmt.Sprintf("timeout error: %v", e.Err)
}

func (e TimeoutError) Unwrap() error {
	return e.Err
}

func MapRequestError(err error) error {
	if err == nil {
		return nil
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return TimeoutError{Err: err}
	}

	return NetworkError{Err: err}
}
