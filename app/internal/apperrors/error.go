package apperrors

import (
	"errors"
	"fmt"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors/codes"
)

type Error struct {
	Code     codes.Code
	Message  error
	Metadata map[string]any
}

func (e *Error) Error() string {
	return e.Message.Error()
}

func (e *Error) GetCode() codes.Code {
	return e.Code
}

func (e *Error) Wrap(msg string) *Error {
	e.Message = fmt.Errorf("%s: %w", msg, e.Message)
	return e
}

//func (e *Error) CheckCode(code codes.Code) *Error {
//	e.Message = fmt.Errorf("%s: %w", msg, e.Message)
//	return e
//}

func (e *Error) Unwrap() error {
	return e.Message
}

func (e *Error) WithMetadata(key string, value any) *Error {
	if e.Metadata == nil {
		e.Metadata = make(map[string]any)
	}

	e.Metadata[key] = value
	return e
}

func New(code codes.Code, message error) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Common predefined errors

func Internal(message error) *Error {
	return New(codes.Internal, message)
}

func InvalidRequest(message error) *Error {
	return New(codes.InvalidRequest, message)
}

func CatNotFound(catID int) *Error {
	return New(codes.CatNotFound, fmt.Errorf("cat with id '%d' was not found", catID))
}

func MissionNotFound(missionID int) *Error {
	return New(codes.MissionNotFound, fmt.Errorf("mission with id '%d' was not found", missionID))
}

func TargetNotFound(target int) *Error {
	return New(codes.TargetNotFound, fmt.Errorf("target with id '%d' was not found", target))
}

func InvalidCatBreed(breed string, reason string) *Error {
	return New(codes.InvalidRequest, fmt.Errorf("cat's breed '%s' is invalid: %s", breed, reason))
}

func InvalidTargetsCount(count, min, max int) *Error {
	return New(codes.InvalidRequest, fmt.Errorf(
		"expected targets count withing range [%d, %d], but got %d",
		min, max, count,
	))
}

func MissionAlreadyCompleted(missionID int) *Error {
	return New(codes.MissionAlreadyCompleted, fmt.Errorf("mission with id '%d' is already completed", missionID))
}

func CatAlreadyAssigned(missionID int) *Error {
	return New(codes.CatAlreadyAssigned, fmt.Errorf("cat with id '%d' is already assigned", missionID))
}

func TargetAlreadyCompleted(targetID int) *Error {
	return New(codes.TargetAlreadyCompleted, fmt.Errorf("target with id '%d' is already completed", targetID))
}

func AllTargetsAreNotCompleted() *Error {
	return New(codes.AllTargetsAreNotCompleted, errors.New("all targets are not completed"))
}
