package api_v1

import (
	"fmt"
	"net/http"
	"sync"
	"time"

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

type getClassPayload struct {
	sync.Mutex
	Name     string                `json:"name"`
	Subject  string                `json:"subject"`
	Lessons  []*lessonClassDetails `json:"lessons"`
	Students []*studentClassDetail `json:"students"`
	Tasks    []*taskClassDetails   `json:"tasks"`
}

type lessonClassDetails struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

type studentClassDetail struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type taskClassDetails struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func GetClassPage(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")

	cl := getClassDetails(uid, cid)

	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	p := new(getClassPayload)
	p.Name = cl.Name
	p.Subject = cl.SubjectName

	convertedStudents := make([]*studentClassDetail, 0)

	for _, cs := range cl.Students {
		convertedStudents = append(convertedStudents, &studentClassDetail{
			Name: fmt.Sprintf("%v %v", cs.FirstName, cs.LastName),
			ID:   cs.ID.String(),
		})
	}

	p.Students = convertedStudents

	var wg = new(sync.WaitGroup)
	wg.Add(2)
	go getTasks(cid, p, wg)
	go getLessons(cid, p, wg)
	wg.Wait()

	return c.JSON(p)
}

func getClassDetails(accID, classID string) *model.Class {
	// Get class details
	var class = new(model.Class)
	// TODO: Handle errors
	res := db.Conn.Preload("Students").Where("account_id = ?", accID).Where("id = ?", classID).First(class)
	if res.RowsAffected == 0 {
		return nil
	}

	return class
}

func getTasks(cid string, p *getClassPayload, wg *sync.WaitGroup) {
	defer wg.Done()
	// Get task details
	var tasks = new([]model.Task)
	// TODO: Handle errors
	db.Conn.Where("class_id = ?", cid).Order("due_date ASC").Find(tasks)

	returnTasks := make([]*taskClassDetails, 0)

	for _, task := range *tasks {
		returnTasks = append(returnTasks, &taskClassDetails{
			ID:        task.ID.String(),
			Name:      task.Name,
			Timestamp: task.DueDate,
		})
	}

	p.Lock()
	p.Tasks = returnTasks
	p.Unlock()
}

func getLessons(cid string, p *getClassPayload, wg *sync.WaitGroup) {
	defer wg.Done()
	// Get task details
	var lessons = new([]model.Lesson)
	// TODO: Handle errors
	db.Conn.Where("class_id = ?", cid).Order("start_time ASC").Find(lessons)

	returnLessons := make([]*lessonClassDetails, 0)

	for _, l := range *lessons {
		returnLessons = append(returnLessons, &lessonClassDetails{
			ID:        l.ID.String(),
			Name:      l.Name,
			Timestamp: l.StartTime,
		})
	}

	p.Lock()
	p.Lessons = returnLessons
	p.Unlock()
}

func DeleteClass(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")

	res := db.Conn.Where("account_id = ?", uid).Where("id = ?", cid).Delete(&model.Class{})
	if res.Error != nil {
		return c.Status(http.StatusBadRequest).SendString("You cannot delete a class while students are still part of it.")
	}

	if res.RowsAffected == 0 {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	return c.Status(http.StatusNoContent).Send(make([]byte, 0))
}

type updateClassPayload struct {
	Name    *string `json:"name,omitempty"`
	Subject *string `json:"subject,omitempty"`
}

func UpdateClass(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")

	p := new(updateClassPayload)
	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if p.Name == nil && p.Subject == nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	cl := new(model.Class)
	res := db.Conn.Where("account_id = ?", uid).Where("id = ?", cid).First(cl)

	if res.RowsAffected == 0 {
		return c.SendStatus(http.StatusNotFound)
	}

	if p.Name != nil {
		cl.Name = *p.Name
		db.Conn.Model(&model.Class{}).Where("id = ?", cl.ID).Update("name", cl.Name)
	}

	if p.Subject != nil {
		cl.SubjectName = *p.Subject
		db.Conn.Model(&model.Class{}).Where("id = ?", cl.ID).Update("subject_name", cl.SubjectName)
	}

	return c.JSON(&dashboardClass{
		ID:      cl.ID.String(),
		Name:    cl.Name,
		Subject: cl.SubjectName,
	})
}
