package middleware

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func newErrBadRequest(message string) error {
	return &fiber.Error{Code: http.StatusBadRequest, Message: message}
}

func newErrNotFound(resource string) error {
	return &fiber.Error{Code: http.StatusNotFound, Message: fmt.Sprintf("%s not found", resource)}
}
