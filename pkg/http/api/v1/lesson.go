package api_v1

import (
	"net/http"
	"time"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type createLessonPayload struct {
	Name        string    `json:"name" validate:"required"`
	Description *string   `json:"description"`
	StartTime   time.Time `json:"startTime" validate:"required"`
	EndTime     time.Time `json:"endTime" validate:"required"`
}

func CreateLesson(c *fiber.Ctx) error {
	cuid, err := uuid.FromString(c.Params("classid"))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	class := new(model.Class)
	ver := db.Conn.Where("id = ?", cuid).First(class)
	if ver.RowsAffected != 1 {
		return c.SendStatus(http.StatusNotFound)
	}

	d := new(createLessonPayload)
	if err := c.BodyParser(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := validate.Struct(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// New Lesson
	nl := &model.Lesson{
		Name:      d.Name,
		StartTime: d.StartTime,
		EndTime:   d.EndTime,
		ClassID:   class.ID,
	}

	if d.Description != nil {
		nl.Description = *d.Description
	}

	// Create the task
	res := db.Conn.Create(nl)

	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(&taskClassDetails{
		ID:        nl.ID.String(),
		Name:      nl.Name,
		Timestamp: nl.StartTime,
	})
}
