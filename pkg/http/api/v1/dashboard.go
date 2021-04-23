package api_v1

import (
	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
)

type dashboardClass struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

type getDashboardReturnPayload struct {
	Classes []*dashboardClass `json:"classes"`
}

// GetDashboardInfo returns the Classman dashboard info
func GetDashboardInfo(c *fiber.Ctx) error {
	uid := c.Locals("uid")

	classes := new([]model.Class)
	db.Conn.Where("account_id = ?", uid).Find(classes)

	returnClasses := make([]*dashboardClass, 0)

	for _, class := range *classes {
		returnClasses = append(returnClasses, &dashboardClass{
			ID:      class.ID.String(),
			Name:    class.Name,
			Subject: class.SubjectName,
		})
	}

	return c.JSON(&getDashboardReturnPayload{
		Classes: returnClasses,
	})
}
