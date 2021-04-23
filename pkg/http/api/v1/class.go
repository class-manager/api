package api_v1

import (
	"net/http"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type createClassData struct {
	Name    string `json:"name" validate:"required"`
	Subject string `json:"subject" validate:"required"`
}

// CreateClass creates a new Classman class
func CreateClass(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	d := new(createClassData)
	if err := c.BodyParser(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := validate.Struct(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// New class
	nc := &model.Class{AccountID: uuid.FromStringOrNil(uid), Name: d.Name, SubjectName: d.Subject}

	res := db.Conn.Create(nc)
	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(&dashboardClass{
		ID:      nc.ID.String(),
		Name:    nc.Name,
		Subject: nc.SubjectName,
	})
}
