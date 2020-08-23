package http

import (
	"github.com/gofiber/fiber"
	uuid "github.com/satori/go.uuid"
)

func idFrom(name string) func(*fiber.Ctx) {
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
