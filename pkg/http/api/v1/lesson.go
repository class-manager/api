package api_v1

import (
	"fmt"
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

type getLessonClassData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type getLessonPayload struct {
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	StartTime time.Time             `json:"startTime"`
	EndTime   time.Time             `json:"endTime"`
	ClassData getLessonClassData    `json:"classData"`
	Students  []*studentClassDetail `json:"students"`
}

// GET /classes/:classid/lessons/:lessonid
func GetLesson(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")
	lid := c.Params("lessonid")

	// Get class data
	cl := getClassDetails(uid, cid)

	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	// Get lesson data
	l := new(model.Lesson)
	ver := db.Conn.Where("id = ?", lid).First(l)
	if ver.RowsAffected != 1 {
		return c.SendStatus(http.StatusNotFound)
	}

	// Create return data payload
	p := new(getLessonPayload)
	p.ClassData = getLessonClassData{
		ID:   cl.ID.String(),
		Name: cl.Name,
	}
	p.EndTime = l.EndTime
	p.ID = l.ID.String()
	p.Name = l.Name
	p.StartTime = l.StartTime

	// Convert students
	convertedStudents := make([]*studentClassDetail, 0)

	for _, cs := range cl.Students {
		convertedStudents = append(convertedStudents, &studentClassDetail{
			Name: fmt.Sprintf("%v %v", cs.FirstName, cs.LastName),
			ID:   cs.ID.String(),
		})
	}

	p.Students = convertedStudents

	return c.JSON(p)
}

// DELETE /classes/:classid/lessons/:lessonid
func DeleteLesson(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")
	lid := c.Params("lessonid")

	// Get class data
	cl := getClassDetails(uid, cid)

	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	res := db.Conn.Where("class_id = ?", cid).Where("id = ?", lid).Delete(&model.Lesson{})
	if res.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	return c.Status(http.StatusNoContent).Send(make([]byte, 0))
}
