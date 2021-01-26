package httpserverutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// HandlerError is an error returned from
// a rest route handler or a middleware.
type HandlerError struct {
	Code          int
	Cause         error
	ClientMessage string
}

func (hErr *HandlerError) Error() string {
	return hErr.Cause.Error()
}

// NewHandlerError returns a HandlerError with the given code and message.
func NewHandlerError(code int, err error) error {
	return &HandlerError{
		Code:          code,
		Cause:         err,
		ClientMessage: err.Error(),
	}
}

// NewHandlerErrorWithCustomClientMessage returns a HandlerError with
// the given code, message and client error message.
func NewHandlerErrorWithCustomClientMessage(code int, err error, clientMessage string) error {
	return &HandlerError{
		Code:          code,
		Cause:         err,
		ClientMessage: clientMessage,
	}
}

// NewInternalServerHandlerError returns a HandlerError with
// the given message, and the http.StatusInternalServerError
// status text as client message.
func NewInternalServerHandlerError(err error) error {
	return NewHandlerErrorWithCustomClientMessage(http.StatusInternalServerError, err, http.StatusText(http.StatusInternalServerError))
}

// NewErrorFromDBErrors takes a slice of database errors and a prefix, and
// returns an error with all of the database errors formatted to one string with
// the given prefix
func NewErrorFromDBErrors(prefix string, dbErrors []error) error {
	dbErrorsStrings := make([]string, len(dbErrors))
	for i, dbErr := range dbErrors {
		dbErrorsStrings[i] = fmt.Sprintf("\"%s\"", dbErr)
	}
	return errors.Errorf("%s [%s]", prefix, strings.Join(dbErrorsStrings, ","))
}

// ClientError is the http response that is sent to the
// client in case of an error.
type ClientError struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func (err *ClientError) Error() string {
	return fmt.Sprintf("%s (Code: %d)", err.ErrorMessage, err.ErrorCode)
}

// SendErr takes a HandlerError and create a ClientError out of it that is sent
// to the http client.
func SendErr(ctx *ServerContext, w http.ResponseWriter, err error) {
	var hErr *HandlerError
	if ok := errors.As(err, &hErr); !ok {
		hErr = NewInternalServerHandlerError(err).(*HandlerError)
	}
	ctx.Warnf("got error: %s", err)
	w.WriteHeader(hErr.Code)
	SendJSONResponse(w, &ClientError{
		ErrorCode:    hErr.Code,
		ErrorMessage: hErr.ClientMessage,
	})
}

// SendJSONResponse encodes the given response to JSON format and
// sends it to the client
func SendJSONResponse(w http.ResponseWriter, response interface{}) {
	b, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprint(w, string(b))
	if err != nil {
		panic(err)
	}
}
