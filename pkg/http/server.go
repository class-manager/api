package http_server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"

	api_v1 "github.com/class-manager/api/pkg/http/api/v1"
	"github.com/class-manager/api/pkg/http/middleware"
)

func Start() {
	app := fiber.New()

	registerMiddleware(app)
	registerRoutes(app)

	// https://github.com/gofiber/recipes/blob/master/graceful-shutdown/main.go
	// Listen from a different goroutine
	go func() {
		if err := app.Listen(":3001"); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)                    // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // When an interrupt or termination signal is sent, notify the channel

	<-c // This blocks the main thread until an interrupt is received
	fmt.Println("Gracefully shutting down...")
	_ = app.Shutdown()

	fmt.Println("Fiber was successful shutdown.")
}

func registerMiddleware(app *fiber.App) {
	// Compression
	app.Use(compress.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
	}))
}

func registerRoutes(app *fiber.App) {
	apiGroup := app.Group("/api")
	registerV1Routes(apiGroup.Group("/v1"))
}

func registerV1Routes(r fiber.Router) {
	// Health endpoint
	r.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	r.Post("/auth/login", api_v1.Login)
	r.Post("/auth/reauth", api_v1.Reauth)
	r.Post("/auth/logout", api_v1.Logout)

	r.Post("/accounts/register", api_v1.Register)

	r.Get("/dashboard", middleware.Protected, api_v1.GetDashboardInfo)

	r.Post("/classes", middleware.Protected, api_v1.CreateClass)
	r.Get("/classes/:classid", middleware.Protected, api_v1.GetClassPage)

	r.Post("/classes/:classid/tasks", middleware.Protected, api_v1.CreateTask)
	r.Get("/classes/:classid/tasks/:taskid", middleware.Protected, api_v1.GetTask)
	r.Patch("/classes/:classid/tasks/:taskid", middleware.Protected, api_v1.UpdateTask)
	r.Delete("/classes/:classid/tasks/:taskid", middleware.Protected, api_v1.DeleteTask)
	r.Patch("/classes/:classid/tasks/:taskid/scores", middleware.Protected, api_v1.UpdateTaskScores)

	r.Patch("/classes/:classid", middleware.Protected, api_v1.UpdateClass)
	r.Delete("/classes/:classid", middleware.Protected, api_v1.DeleteClass)

	r.Get("/classes/:classid/students", middleware.Protected, api_v1.GetStudentsFromClass)
	r.Post("/classes/:classid/students", middleware.Protected, api_v1.AddStudentsToClass)
	r.Delete("/classes/:classid/students", middleware.Protected, api_v1.DeleteStudentsFromClass)

	r.Post("/classes/:classid/lessons", middleware.Protected, api_v1.CreateLesson)
	r.Get("/classes/:classid/lessons/:lessonid", middleware.Protected, api_v1.GetLesson)
	r.Patch("/classes/:classid/lessons/:lessonid", middleware.Protected, api_v1.UpdateLesson)
	r.Delete("/classes/:classid/lessons/:lessonid", middleware.Protected, api_v1.DeleteLesson)

	r.Get("/classes/:classid/lessons/:lessonid/students/:studentid", middleware.Protected, api_v1.GetStudentForLesson)
	r.Patch("/classes/:classid/lessons/:lessonid/students/:studentid", middleware.Protected, api_v1.UpdateStudentForLesson)

	r.Get("/students", middleware.Protected, api_v1.GetStudents)
	r.Get("/students/:studentid", middleware.Protected, api_v1.GetStudent)
	r.Patch("/students/:studentid", middleware.Protected, api_v1.UpdateStudent)
	r.Delete("/students/:studentid", middleware.Protected, api_v1.DeleteStudent)
	r.Post("/students", middleware.Protected, api_v1.CreateStudent)
}
