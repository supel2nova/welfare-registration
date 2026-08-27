package httpx

import "github.com/supel2nova/welfare-registration/backend/pkg/apperror"

const (
	statusSuccess = "0"
	statusError   = "1"
)

type Envelope struct {
	Data              any                   `json:"data"`
	StatusCode        string                `json:"statusCode"`
	StatusDescription string                `json:"statusDescription"`
	ErrorCode         string                `json:"errorCode,omitempty"`
	ErrorMessage      string                `json:"errorMessage,omitempty"`
	FieldErrors       []apperror.FieldError `json:"fieldErrors,omitempty"`
}

func OK(data any) Envelope {
	return Envelope{Data: data, StatusCode: statusSuccess, StatusDescription: "Success"}
}

func Fail(e *apperror.Error) (int, Envelope) {
	return e.HTTPStatus, Envelope{
		Data:              e.Data,
		StatusCode:        statusError,
		StatusDescription: e.Message,
		ErrorCode:         e.Code,
		ErrorMessage:      e.Message,
		FieldErrors:       e.FieldErrors,
	}
}
