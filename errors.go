package rpr

import "errors"

// ErrNotAResponse is returned when a response validation fails or
// the provided object does not comply with the expected Response structure.
var ErrNotAResponse = errors.New("could not validate Response")

// ErrUnpackPayloadFail is returned when the internal type assignment
// or data unpacking mechanism is unable to extract the underlying payload structure.
var ErrUnpackPayloadFail = errors.New("unable do unpack Payload")
