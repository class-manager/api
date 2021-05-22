package api_v1

import (
	"net/http"
	"time"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type createStudentPayload struct {
	FirstName       string    `json:"firstName" validate:"required"`
	LastName        string    `json:"lastName" validate:"required"`
	DOB             time.Time `json:"dob" validate:"required"`
	GraduatingClass uint32    `json:"graduatingClass" validate:"required"`
	GeneralNote     *string   `json:"generalNote"`
	StudentNumber   *string   `json:"studentNumber"`
}

// POST /api/v1/students
func CreateStudent(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)

	s := new(createStudentPayload)
	if err := c.BodyParser(s); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := validate.Struct(s); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	ns := &model.Student{
		FirstName:       s.FirstName,
		LastName:        s.LastName,
		DOB:             s.DOB,
		GraduatingClass: s.GraduatingClass,
		CreatedByID:     uuid.FromStringOrNil(uid),
	}

	if s.GeneralNote != nil {
		ns.GeneralNote = *s.GeneralNote
	}

	if s.StudentNumber != nil {
		ns.StudentNumber = *s.StudentNumber
	}

	// Create the task
	res := db.Conn.Create(ns)

	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(&createStudentPayload{
		FirstName:       s.FirstName,
		LastName:        s.LastName,
		DOB:             s.DOB,
		GraduatingClass: s.GraduatingClass,
		GeneralNote:     s.GeneralNote,
		StudentNumber:   s.StudentNumber,
	})
}
