package httpfiber

import (
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
)

const requestIDHeader = "X-Request-Id"

func requestIDMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		requestID := c.Get(requestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}
		c.Locals("request_id", requestID)
		c.Set(requestIDHeader, requestID)
		return c.Next()
	}
}

func requestLogMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		started := time.Now()
		err := c.Next()
		requestID, _ := c.Locals("request_id").(string)
		log.Printf("http method=%s path=%s status=%d latency=%s request_id=%s", c.Method(), c.Path(), c.Response().StatusCode(), time.Since(started), requestID)
		return err
	}
}

func authMiddleware(auth ports.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Missing bearer token")
		}

		token := strings.TrimPrefix(header, "Bearer ")
		requestID := c.Locals("request_id").(string)
		ctx, err := auth.GetMe(c.UserContext(), token, requestID)
		if err != nil {
			return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		}

		c.Locals("auth_context", ctx)
		return c.Next()
	}
}
