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

type createTaskPayload struct {
	Name        string     `json:"name" validate:"required"`
	Description *string    `json:"description"`
	OpenDate    *time.Time `json:"openDate"`
	DueDate     time.Time  `json:"dueDate" validate:"required"`
	MaxMark     int32      `json:"maxMark" validate:"required"`
}

type returnTask struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Class string `json:"class"`
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
		Name:     d.Name,
		Open:     true,
		DueDate:  d.DueDate,
		MaxMark:  d.MaxMark,
		ClassID:  class.ID,
		OpenDate: time.Now(),
	}

	if d.Description != nil {
		nt.Description = d.Description
	}

	if d.OpenDate != nil {
		nt.OpenDate = *d.OpenDate
	}

	// Create the task
	res := db.Conn.Create(nt)

	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.Status(http.StatusCreated).JSON(&returnTask{
		ID:    nt.ID.String(),
		Name:  nt.Name,
		Class: class.Name,
	})
}

type studentResultData struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Score *float64 `json:"score"`
}

type taskPageReturnData struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	Description    *string             `json:"description"`
	OpenDate       time.Time           `json:"openDate"`
	DueDate        time.Time           `json:"dueDate"`
	MaxMark        int32               `json:"maxMark"`
	StudentResults []studentResultData `json:"studentResults"`
	ClassName      string              `json:"className"`
	ClassID        string              `json:"classID"`
}

func GetTask(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")
	tid := c.Params("taskid")

	cl := getClassDetails(uid, cid)

	// Get class and ensure it exists
	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	// Ensure task exists
	t := new(model.Task)

	res := db.Conn.Where("id = ?", tid).Where("class_id = ?", cid).First(t)
	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	if t.ID == uuid.Nil {
		return c.SendStatus(http.StatusNotFound)
	}

	// Get all task results for this task
	tr := new([]model.TaskResult)
	res = db.Conn.Where("task_id = ?", tid).Find(tr)
	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Get map of task results
	resultsMap := make(map[string]float64)
	for _, r := range *tr {
		resultsMap[r.ID.String()] = r.Mark
	}

	// Create list of student results to return
	sr := make([]studentResultData, 0)
	for _, s := range cl.Students {
		r := studentResultData{
			ID:    s.ID.String(),
			Score: nil,
			Name:  fmt.Sprintf("%v %v", s.FirstName, s.LastName),
		}

		if val, ok := resultsMap[s.ID.String()]; ok {
			r.Score = &val
		}

		sr = append(sr, r)
	}

	// Return data to user
	rd := &taskPageReturnData{
		ID:             t.ID.String(),
		Name:           t.Name,
		Description:    nil,
		OpenDate:       t.OpenDate,
		DueDate:        t.DueDate,
		MaxMark:        t.MaxMark,
		StudentResults: sr,
		ClassName:      cl.Name,
		ClassID:        cl.ID.String(),
	}

	if t.Description != nil {
		rd.Description = t.Description
	}

	return c.JSON(rd)
}

// PATCH /api/v1/classes/:classid/tasks/:taskid
func UpdateTask(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	cid := c.Params("classid")
	tid := c.Params("taskid")

	cl := getClassDetails(uid, cid)

	// Get class and ensure it exists
	if cl == nil {
		return c.Status(http.StatusNotFound).Send(make([]byte, 0))
	}

	// Ensure task exists
	t := new(model.Task)

	res := db.Conn.Where("id = ?", tid).Where("class_id = ?", cid).First(t)
	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	if t.ID == uuid.Nil {
		return c.SendStatus(http.StatusNotFound)
	}

	d := new(taskPageReturnData)
	if err := c.BodyParser(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := validate.Struct(d); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	t.Name = d.Name
	t.Description = d.Description
	t.DueDate = d.DueDate
	t.MaxMark = d.MaxMark

	res = db.Conn.Save(t)
	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Get all task results for this task
	tr := new([]model.TaskResult)
	res = db.Conn.Where("task_id = ?", tid).Find(tr)
	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Get map of task results
	resultsMap := make(map[string]float64)
	for _, r := range *tr {
		resultsMap[r.ID.String()] = r.Mark
	}

	// Create list of student results to return
	sr := make([]studentResultData, 0)
	for _, s := range cl.Students {
		r := studentResultData{
			ID:    s.ID.String(),
			Score: nil,
			Name:  fmt.Sprintf("%v %v", s.FirstName, s.LastName),
		}

		if val, ok := resultsMap[s.ID.String()]; ok {
			r.Score = &val
		}

		sr = append(sr, r)
	}

	// Return data to user
	rd := &taskPageReturnData{
		ID:             t.ID.String(),
		Name:           t.Name,
		Description:    nil,
		OpenDate:       t.OpenDate,
		DueDate:        t.DueDate,
		MaxMark:        t.MaxMark,
		StudentResults: sr,
		ClassName:      cl.Name,
		ClassID:        cl.ID.String(),
	}

	return c.JSON(rd)
}
