package rpr

import "errors"

var ErrNotAResponse = errors.New("could not validate Response")
var ErrUnpackPayloadFail = errors.New("unable do unpack Payload")
