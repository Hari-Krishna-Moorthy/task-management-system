package task_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/controller"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	test_helpers "github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/test"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Task", func() {
	var (
		app              *fiber.App
		createTaskParams *types.CreateTaskRequest
		createTaskBody   []byte
		token            string
	)

	BeforeEach(func() {
		app = fiber.New()
		app.Post("/task", controller.GetTaskController(services.GetTaskService(services.GetTaskRepository(test_helpers.GetTestDatabase()))).CreateTask)

		createTaskParams = &types.CreateTaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     "2025-01-01",
		}

		bytes, _ := bcrypt.GenerateFromPassword([]byte("validPassword123"), bcrypt.DefaultCost)

		services.GetAuthRepository(test_helpers.GetTestDatabase()).CreateUser(context.TODO(), &models.User{
			ID:        helpers.GenerateUUID(),
			Username:  "validUsername",
			Email:     "validemail@example.com",
			Password:  string(bytes),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})

		user, _ := services.GetAuthRepository(test_helpers.GetTestDatabase()).FindUserByUsername(context.TODO(), "validUsername")
		token, _ = test_helpers.GenerateToken(user)

	})

	JustBeforeEach(func() {
		createTaskBody, _ = json.Marshal(createTaskParams)
	})

	When("request is valid", func() {
		It("should create a new task", func() {
			req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(createTaskBody))
			req.Header.Set("Content-Type", "application/json")
			cookie := &http.Cookie{
				Name:     utils.CookieKeyToken,
				Value:    token,
				Expires:  time.Now().Add(utils.JWT_TOKEN_EXPIRY),
				HttpOnly: true,
				Domain:   config.GetConfig().Server.AppDomain,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
				MaxAge:   int(utils.JWT_TOKEN_EXPIRY.Seconds()),
			}
			req.AddCookie(cookie)

			res, err := app.Test(req)

			log.Printf("res: %v, err: %v\n", res, err)

			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(fiber.StatusOK))

			var response types.CreateTaskResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response.Success).To(BeTrue())
			Expect(response.Message).To(Equal("Task created successfully"))

			// Assert the task is created in the database

			task := &models.Task{}
			taskCollection := test_helpers.GetTestDatabase().Collection("tasks")
			taskCollection.FindOne(context.TODO(), bson.M{"title": createTaskParams.Title}).Decode(&task)

			Expect(task.Title).To(Equal(createTaskParams.Title))
			Expect(task.Description).To(Equal(createTaskParams.Description))
			Expect(task.DueDate.Format(utils.TimeLayout)).To(Equal(createTaskParams.DueDate))
		})
	})

	When("request is invalid", func() {
		It("Not authorized", func() {
			req, err := http.NewRequest("POST", "/task", bytes.NewBuffer(createTaskBody))
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(fiber.StatusNotFound))
		})
	})

})
