package httpfiber

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/isapr/mini-erp/services/api-gateway/internal/domain"
	"github.com/isapr/mini-erp/services/api-gateway/internal/ports"
	"github.com/isapr/mini-erp/services/api-gateway/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	auth ports.AuthService
}

func NewAuthHandler(auth ports.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Signup(c *fiber.Ctx) error {
	var req domain.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	requestID := c.Locals("request_id").(string)
	session, err := h.auth.Signup(c.UserContext(), req, requestID)
	if err != nil {
		return mapAuthError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(session)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "Invalid JSON body")
	}
	req.Email = strings.TrimSpace(req.Email)
	requestID := c.Locals("request_id").(string)
	session, err := h.auth.Login(c.UserContext(), req, requestID)
	if err != nil {
		return mapAuthError(c, err)
	}
	return c.JSON(session)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Get("X-Refresh-Token")
	if refreshToken == "" {
		return writeError(c, fiber.StatusBadRequest, "MISSING_REFRESH_TOKEN", "Missing refresh token")
	}
	requestID := c.Locals("request_id").(string)
	session, err := h.auth.Refresh(c.UserContext(), refreshToken, requestID)
	if err != nil {
		return mapAuthError(c, err)
	}
	return c.JSON(session)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	header := c.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return writeError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Missing bearer token")
	}
	if err := h.auth.Logout(c.UserContext(), strings.TrimPrefix(header, "Bearer ")); err != nil {
		return writeError(c, fiber.StatusInternalServerError, "LOGOUT_FAILED", err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	ctx := c.Locals("auth_context")
	return c.JSON(ctx)
}

func mapAuthError(c *fiber.Ctx, err error) error {
	switch {
	case status.Code(err) == codes.AlreadyExists && strings.Contains(err.Error(), "business"):
		return writeError(c, fiber.StatusConflict, "BUSINESS_ALREADY_EXISTS", "Business already exists")
	case status.Code(err) == codes.AlreadyExists && strings.Contains(err.Error(), "user"):
		return writeError(c, fiber.StatusConflict, "USER_ALREADY_EXISTS", "User already exists")
	case status.Code(err) == codes.InvalidArgument:
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case status.Code(err) == codes.Unauthenticated:
		return writeError(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials")
	case errors.Is(err, service.ErrValidation):
		return writeError(c, fiber.StatusBadRequest, "VALIDATION_ERROR", "Validation failed")
	case errors.Is(err, service.ErrInvalidCredentials):
		return writeError(c, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials")
	case errors.Is(err, service.ErrInvalidToken):
		return writeError(c, fiber.StatusUnauthorized, "INVALID_TOKEN", "Invalid token")
	default:
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
