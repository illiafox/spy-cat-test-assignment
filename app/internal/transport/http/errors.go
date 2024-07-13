package http

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors/codes"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var errorCodesToHTTP = map[codes.Code]int{
	codes.Internal:                  http.StatusInternalServerError,
	codes.InvalidRequest:            http.StatusBadRequest,
	codes.CatNotFound:               http.StatusNotFound,
	codes.MissionNotFound:           http.StatusNotFound,
	codes.TargetNotFound:            http.StatusNotFound,
	codes.MissionAlreadyCompleted:   http.StatusForbidden,
	codes.CatAlreadyAssigned:        http.StatusForbidden,
	codes.AllTargetsAreNotCompleted: http.StatusForbidden,
	codes.TargetAlreadyCompleted:    http.StatusForbidden,
}

type Error struct {
	Ok       bool           `json:"ok"`
	Code     codes.Code     `json:"code"`
	Message  string         `json:"message"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

func extractError(err error) *apperrors.Error {
	var e *apperrors.Error
	errors.As(err, &e)
	return e
}

func toHTTPError(ctx *fiber.Ctx, e *apperrors.Error) error {
	status, ok := errorCodesToHTTP[e.Code]
	if !ok {
		return fmt.Errorf("error '%s': code '%s' is not in errorCodesToHTTP", e.Message, e.Code)
	}

	respError := Error{
		Ok:       false,
		Code:     e.Code,
		Message:  e.Message.Error(),
		Metadata: e.Metadata,
	}

	return ctx.Status(status).JSON(respError)
}

func RespondWithError(ctx *fiber.Ctx, err error) error {
	serviceError := extractError(err)
	if serviceError == nil { // process as internal
		return err
	}

	if serviceError.Code == codes.Internal {
		return err // pass to error middleware
	}

	return toHTTPError(ctx, serviceError)
}

func ErrorMiddleware(logger *zap.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err != nil {
			if errors.As(err, new(*fiber.Error)) {
				return err
			}

			fields := []zap.Field{
				zap.Error(err),
				zap.String("url", strings.Clone(ctx.OriginalURL())),
			}

			var serviceError *apperrors.Error
			if errors.As(err, &serviceError) {
				if serviceError.Metadata != nil {
					for k, v := range serviceError.Metadata {
						fields = append(fields, zap.Any(k, v))
					}
				}

			}

			logger.Error("failed to process request", fields...)

			internalError := Error{
				Ok:       false,
				Code:     codes.Internal,
				Message:  "internal error, please try again later",
				Metadata: nil,
			}

			if sErr := ctx.Status(http.StatusInternalServerError).JSON(internalError); sErr != nil {
				logger.Error("Failed to send internal error message", zap.Error(err))
			}

			return nil
		}

		return nil
	}
}
