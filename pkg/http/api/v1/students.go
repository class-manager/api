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
		ns.GeneralNote = s.GeneralNote
	}

	if s.StudentNumber != nil {
		ns.StudentNumber = s.StudentNumber
	}

	// Create the task
	res := db.Conn.Create(ns)

	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(&createStudentPayload{
		ID:              ns.ID.String(),
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
			GeneralNote:     s.GeneralNote,
			StudentNumber:   s.StudentNumber,
		})
	}

	return c.JSON(returnStudents)
}

type studentIDList struct {
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

	ss := new(studentIDList)
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

	return c.SendStatus(http.StatusOK)
}

// DELETE /api/v1/classes/:classid/students
func DeleteStudentsFromClass(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")

	cl := getClassDetails(uid, cid)

	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	ss := new(studentIDList)
	if err := c.BodyParser(ss); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := validate.Struct(ss); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	tx := db.Conn.Begin()
	for _, id := range ss.Students {
		tx.Exec("DELETE FROM students_classes WHERE student_id = ? AND class_id = ?", id, cid)
	}

	res := tx.Commit()
	if res.Error != nil {
		tx.Rollback()
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

// GET /api/v1/classes/:classid/students
func GetStudentsFromClass(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")

	cl := getClassDetails(uid, cid)

	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	d := make([]*studentClassDetail, 0)

	for _, s := range cl.Students {
		d = append(d, &studentClassDetail{
			ID:   s.ID.String(),
			Name: fmt.Sprintf("%v %v", s.FirstName, s.LastName),
		})
	}

	return c.JSON(d)
}

type studentLessonPayload struct {
	ID              string    `json:"id"`
	BHID            *string   `json:"bhid"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	DOB             time.Time `json:"dob"`
	GeneralNote     *string   `json:"generalNote"`
	BehaviouralNote *string   `json:"behaviouralNote"`
}

// GET Protected::/classes/:classid/lessons/:lessonid/student/:studentid
func GetStudentForLesson(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	lid := c.Params("lessonid")
	sid := c.Params("studentid")

	// Check if the student exists
	s := new(model.Student)
	db.Conn.Where("id = ?", sid).Where("created_by_id = ?", uid).First(s)

	if s.ID == uuid.Nil {
		return c.SendStatus(http.StatusNotFound)
	}

	// Student exists, get notes for this student
	bn := new(model.BehaviourNote)
	db.Conn.Where("lesson_id = ?", lid).Where("student_id = ?", sid).Find(bn)

	// Return data
	rd := studentLessonPayload{
		GeneralNote:     s.GeneralNote,
		ID:              s.ID.String(),
		BHID:            nil,
		FirstName:       s.FirstName,
		LastName:        s.LastName,
		DOB:             s.DOB,
		BehaviouralNote: nil,
	}

	if bn.ID != uuid.Nil {
		bhid := bn.ID.String()
		rd.BHID = &bhid
		rd.BehaviouralNote = &bn.Note
	}

	return c.JSON(rd)
}

type updateStudentForLessonDetails struct {
	BehaviourNote *string `json:"behaviouralNote"`
}

// PATCH Protected::/classes/:classid/lessons/:lessonid/student/:studentid
func UpdateStudentForLesson(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	lid := c.Params("lessonid")
	sid := c.Params("studentid")

	// Parse new details
	pd := new(updateStudentForLessonDetails)
	if err := c.BodyParser(pd); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// Check if the student exists
	s := new(model.Student)
	db.Conn.Where("id = ?", sid).Where("created_by_id = ?", uid).First(s)

	if s.ID == uuid.Nil {
		return c.SendStatus(http.StatusNotFound)
	}

	// Student exists, get notes for this student
	bn := new(model.BehaviourNote)
	db.Conn.Where("lesson_id = ?", lid).Where("student_id = ?", sid).Find(bn)

	// Update data
	if pd.BehaviourNote != nil {
		if bn.ID != uuid.Nil {
			bn.Note = *pd.BehaviourNote

			res := db.Conn.Save(bn)
			if res.Error != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}
		} else {
			bn = &model.BehaviourNote{
				Note:      *pd.BehaviourNote,
				LessonID:  uuid.FromStringOrNil(lid),
				StudentID: s.ID,
			}
			res := db.Conn.Create(bn)

			if res.Error != nil {
				return c.SendStatus(http.StatusInternalServerError)
			}
		}
		// Delete the behaviour note if it exists
	} else if bn.ID != uuid.Nil {
		db.Conn.Delete(bn)
	}

	// Return data
	rd := studentLessonPayload{
		GeneralNote:     s.GeneralNote,
		ID:              s.ID.String(),
		BHID:            nil,
		FirstName:       s.FirstName,
		LastName:        s.LastName,
		DOB:             s.DOB,
		BehaviouralNote: nil,
	}

	if bn.ID != uuid.Nil {
		bhid := bn.ID.String()
		rd.BHID = &bhid
		rd.BehaviouralNote = &bn.Note
	}

	return c.JSON(rd)
}
