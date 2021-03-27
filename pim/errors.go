package pim

import "fmt"

// ErrFailed is a general error that signifies an error event,
// for which there's no other defined error. Acts as a default error.
const ErrFailed = clientError("PIMClient failed")

type clientError string

func (err clientError) Error() string {
	return string(err)
}

func wrapFailedError() error {
	return fmt.Errorf("%w", ErrFailed)
}
