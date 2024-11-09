package task_test

import (
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
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get Tasks", func() {
	var (
		app              *fiber.App
		token            string
		user             *models.User
		task1ID, task2ID string
	)

	BeforeEach(func() {
		app = fiber.New()
		app.Get("/task/:id", controller.GetTaskController(services.GetTaskService(services.GetTaskRepository(test_helpers.GetTestDatabase()))).GetTask)

		bytes, _ := bcrypt.GenerateFromPassword([]byte("validPassword123"), bcrypt.DefaultCost)

		services.GetAuthRepository(test_helpers.GetTestDatabase()).CreateUser(context.TODO(), &models.User{
			ID:        helpers.GenerateUUID(),
			Username:  "validUsername",
			Email:     "validemail@example.com",
			Password:  string(bytes),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		})

		user, _ = services.GetAuthRepository(test_helpers.GetTestDatabase()).FindUserByUsername(context.TODO(), "validUsername")
		token, _ = test_helpers.GenerateToken(user)
		task1ID, task2ID = helpers.GenerateUUID(), helpers.GenerateUUID()
	})

	Context("When getting tasks", func() {
		BeforeEach(func() {
			services.GetTaskRepository(test_helpers.GetTestDatabase()).CreateTask(context.TODO(), &models.Task{
				ID:          task1ID,
				Title:       "Test Task",
				Description: "Test Description",
				DueDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				UserID:      user.ID,
			})
			services.GetTaskRepository(test_helpers.GetTestDatabase()).CreateTask(context.TODO(), &models.Task{
				ID:          task2ID,
				Title:       "Test Task 2",
				Description: "Test Description 2",
				DueDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				UserID:      user.ID,
			})
		})
		It("should return a task", func() {

			uri := "/task/" + task1ID
			req, err := http.NewRequest("GET", uri, nil)

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

			var response types.GetTaskResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response.Success).To(BeTrue())
			Expect(response.Message).To(Equal("Task fetched successfully"))
			Expect(response.Task.ID).To(Equal(task1ID))
			Expect(response.Task.Title).To(Equal("Test Task"))
			Expect(response.Task.Description).To(Equal("Test Description"))
			Expect(response.Task.DueDate).To(Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
			Expect(response.Task.UserID).To(Equal(user.ID))

		})
	})

	When("Task does not exist", func() {
		var taskId string
		BeforeEach(func() {
			taskId = helpers.GenerateUUID()
		})
		It("should return error", func() {
			req, err := http.NewRequest("GET", "/task/"+taskId, nil)
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
			Expect(res.StatusCode).To(Equal(fiber.StatusNotFound))
		})
	})

	When("request is invalid", func() {
		It("Not authorized", func() {
			req, err := http.NewRequest("GET", "/task/"+task2ID, nil)
			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)

			Expect(err).To(BeNil())
			Expect(res.StatusCode).To(Equal(fiber.StatusNotFound))
		})
	})

})
