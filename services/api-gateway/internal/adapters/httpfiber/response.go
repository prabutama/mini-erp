package httpfiber

import "github.com/gofiber/fiber/v2"

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(c *fiber.Ctx, status int, code string, message string) error {
	return c.Status(status).JSON(errorResponse{Code: code, Message: message})
}
