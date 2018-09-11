package pusherplatform

import (
	"fmt"
)

// BodyNotJSONError indicates a response to an elements request has a body that
// could not be parsed as JSON.
type BodyNotJSONError struct {
	JSONDecodeError error
	StatusCode      int
	BodyBytes       []byte
}

func (e BodyNotJSONError) Error() string {
	return fmt.Sprintf(
		"Body is not valid JSON. Status: %v Body: %v Error: %s",
		e.StatusCode,
		e.BodyBytes,
		e.JSONDecodeError,
	)
}
