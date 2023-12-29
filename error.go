package base

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CodeResponse int

const (
	ErrorCodeBadRequest   CodeResponse = 400
	ErrorCodeUnauthorized CodeResponse = 401
	ErrorCodeNotFound     CodeResponse = 404
	ErrorCodeSystemError  CodeResponse = 500
	ErrorCodeForbidden    CodeResponse = 403
)

type Source string

type Error struct {
	Code       CodeResponse `json:"code"`
	Message    string       `json:"message"`
	TraceID    string       `json:"trace_id,omitempty"`
	Detail     string       `json:"detail"`
	Source     Source       `json:"source"`
	HTTPStatus int          `json:"http_status"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:[%d], message:[%s], detail:[%s], source:[%s]", e.Code, e.Message, e.Detail, e.Source)
}

func (e *Error) GetHttpStatus() int {
	return e.HTTPStatus
}

func (e *Error) GetCode() CodeResponse {
	return e.Code
}

func (e *Error) GetMessage() string {
	return e.Message
}

func (e *Error) SetTraceId(traceId string) *Error {
	e.TraceID = fmt.Sprintf("%s:%d", traceId, time.Now().Unix())
	return e
}

func (e *Error) SetHTTPStatus(status int) *Error {
	e.HTTPStatus = status
	return e
}

func (e *Error) SetMessage(msg string) *Error {
	e.Message = msg
	return e
}

func (e *Error) SetDetail(detail string) *Error {
	e.Detail = detail
	return e
}

func (e *Error) GetDetail() string {
	return e.Detail
}

func (e *Error) SetSource(source Source) *Error {
	e.Source = source
	return e
}

func (e *Error) ToJSon() string {
	data, _ := json.Marshal(e)
	return string(data)
}

var (
	// Status 4xx ********

	ErrUnauthorized = func(ctx context.Context) *Error {
		return &Error{
			Code:       ErrorCodeUnauthorized,
			Message:    DefaultUnauthorizedMessage,
			TraceID:    GetTraceId(ctx),
			HTTPStatus: http.StatusUnauthorized,
		}
	}

	ErrNotFound = func(ctx context.Context, object string) *Error {
		if object == "" {
			object = "object"
		}
		return &Error{
			Code:       ErrorCodeNotFound,
			Message:    fmt.Sprintf("%s %s", object, "not found"),
			TraceID:    GetTraceId(ctx),
			HTTPStatus: http.StatusNotFound,
		}
	}

	ErrBadRequest = func(ctx context.Context) *Error {
		return &Error{
			Code:       ErrorCodeBadRequest,
			Message:    DefaultBadRequestMessage,
			TraceID:    GetTraceId(ctx),
			HTTPStatus: http.StatusBadRequest,
		}
	}

	// Status 5xx *******

	ErrSystemError = func(ctx context.Context, detail string) *Error {
		return &Error{
			Code:       ErrorCodeSystemError,
			Message:    DefaultServerErrorMessage,
			TraceID:    GetTraceId(ctx),
			HTTPStatus: http.StatusInternalServerError,
			Detail:     detail,
		}
	}

	ErrForbidden = func(ctx context.Context) *Error {
		return &Error{
			Code:       ErrorCodeForbidden,
			Message:    DefaultForbiddenMessage,
			TraceID:    GetTraceId(ctx),
			HTTPStatus: http.StatusForbidden,
		}
	}
)

const (
	DefaultServerErrorMessage  = "Something has gone wrong, please contact admin"
	DefaultBadRequestMessage   = "Invalid request"
	DefaultUnauthorizedMessage = "Token invalid"
	DefaultForbiddenMessage    = "Forbidden"
)
