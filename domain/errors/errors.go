package errors

import (
	"fmt"
	"strconv"
)

// Error represent domain's error
type Error struct {
	Message string
	Code    string
}

func (e Error) Error() string {
	return e.Message
}

var registeredCodes = make(map[string]struct{})

// NewError creates an error
func NewError(code, msg string) error {
	_, existed := registeredCodes[code]
	if existed {
		panic(fmt.Sprintf("code %s existed", code))
	}

	if len(code) != 5 {
		panic("invalid code")
	}

	_, err := strconv.ParseInt(code, 10, 32)
	if err != nil {
		panic("invalid code")
	}

	rpcCode := code[:2]
	num, err := strconv.ParseInt(rpcCode, 10, 32)
	if err != nil {
		panic("invalid code")
	}
	if num < 1 || num > 16 {
		panic("invalid code")
	}

	registeredCodes[code] = struct{}{}

	return Error{
		Code:    code,
		Message: msg,
	}
}
