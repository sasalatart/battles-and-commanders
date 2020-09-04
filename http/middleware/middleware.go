package middleware

import (
	"github.com/gofiber/fiber"
	uuid "github.com/satori/go.uuid"
)

// IDFrom is a fiber middleware that reads, parses and validates a named UUID in the route params
// and stores it in Locals with the same name
func IDFrom(name string) func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		id, err := uuid.FromString(c.Params(name))
		if err != nil {
			c.Next(fiber.ErrBadRequest)
			return
		}
		c.Locals(name, id)
		c.Next()
	}
}
