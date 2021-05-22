package api_v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

type createStudentPayload struct {
	ID              string    `json:"id"`
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

// GET /api/v1/students
func GetStudents(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	term := c.Query("name")

	var students = new([]model.Student)

	db.Conn.Where("created_by_id = ?", uid).Where(db.Conn.Where("first_name LIKE ?", fmt.Sprintf("%%%v%%", term)).Or("last_name LIKE ?", fmt.Sprintf("%%%v%%", term))).Find(students)

	returnStudents := make([]*createStudentPayload, 0)

	for _, s := range *students {
		returnStudents = append(returnStudents, &createStudentPayload{
			ID:              s.ID.String(),
			FirstName:       s.FirstName,
			LastName:        s.LastName,
			DOB:             s.DOB,
			GraduatingClass: s.GraduatingClass,
			GeneralNote:     &s.GeneralNote,
			StudentNumber:   &s.StudentNumber,
		})
	}

	return c.JSON(returnStudents)
}

type addStudentsToClassPayload struct {
	Students []string `json:"students" validate:"required,unique"`
}

// POST /api/v1/classes/:classid/students
func AddStudentsToClass(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")

	cl := getClassDetails(uid, cid)

	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	ss := new(addStudentsToClassPayload)
	if err := c.BodyParser(ss); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := validate.Struct(ss); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// Find all students in the list of ids provided
	var students = new([]model.Student)
	db.Conn.Preload("Classes").Where("created_by_id = ?", uid).Where("id IN ?", ss.Students).Find(students)

	// Remove all students from list who are already part of the class
	nss := make([]model.Student, 0)
	for _, s := range *students {
		add := true
		for _, c := range s.Classes {
			if c.ID == cl.ID {
				add = false
			}
		}

		if add {
			nss = append(nss, s)
		}
	}

	// Add the students to the set class
	for _, s := range nss {
		s.Classes = append(s.Classes, cl)
	}

	tx := db.Conn.Begin()
	for _, s := range nss {
		tx.Exec("INSERT INTO students_classes VALUES (?, ?)", s.ID, cl.ID)
	}

	res := tx.Commit()
	if res.Error != nil {
		tx.Rollback()
		return c.SendStatus(http.StatusInternalServerError)
	}

	log.Printf("students: %#+v\n", len(*students))
	return c.SendStatus(http.StatusOK)
}
