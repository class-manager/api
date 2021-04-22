package api_v1

import (
	"github.com/class-manager/api/pkg/util/security"
	"github.com/gofiber/fiber/v2"
)

func Logout(c *fiber.Ctx) error {
	security.ClearRefreshCookie(c)

	return nil
}
