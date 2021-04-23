package api_v1

import (
	"net/http"
	"time"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type createTaskPayload struct {
	Name        string     `json:"name" validate:"required"`
	Description *string    `json:"description"`
	OpenDate    *time.Time `json:"openDate"`
	DueDate     time.Time  `json:"dueDate" validate:"required"`
	MaxMark     int32      `json:"maxMark" validate:"required"`
}

func CreateTask(c *fiber.Ctx) error {
	cuid, err := uuid.FromString(c.Params("classid"))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	class := new(model.Class)
	ver := db.Conn.Where("id = ?", cuid).First(class)
	if ver.RowsAffected != 1 {
		return c.SendStatus(http.StatusNotFound)
	}

	d := new(createTaskPayload)
	if err := c.BodyParser(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := validate.Struct(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// New Task
	nt := &model.Task{
		Name:    d.Name,
		Open:    false,
		DueDate: d.DueDate,
		MaxMark: d.MaxMark,
		ClassID: class.ID,
	}

	if d.Description != nil {
		nt.Description = *d.Description
	}

	if d.OpenDate != nil {
		nt.OpenDate = *d.OpenDate
	}

	// Create the task
	res := db.Conn.Create(nt)

	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(&dashboardTask{
		ID:    nt.ID.String(),
		Name:  nt.Name,
		Class: class.Name,
	})
}
