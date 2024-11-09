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
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Update Task", func() {

	var (
		app              *fiber.App
		user             *models.User
		UpdateTaskParams *types.UpdateTaskRequest
		token            string
		UpdateTaskBody   []byte
		taskId           string
	)

	BeforeEach(func() {
		app = fiber.New()
		app.Put("/task/:id", controller.GetTaskController(services.GetTaskService(services.GetTaskRepository(test_helpers.GetTestDatabase()))).UpdateTask)

		UpdateTaskParams = &types.UpdateTaskRequest{
			Title:       "Updated Task",
			Description: "Updated Description",
			DueDate:     "2025-01-01",
			Status:      "Done",
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

		user, _ = services.GetAuthRepository(test_helpers.GetTestDatabase()).FindUserByUsername(context.TODO(), "validUsername")
		token, _ = test_helpers.GenerateToken(user)
		taskId = helpers.GenerateUUID()
		services.GetTaskRepository(test_helpers.GetTestDatabase()).CreateTask(context.TODO(), &models.Task{
			ID:          taskId,
			Title:       "Test Task",
			Description: "Test Description",
			DueDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			UserID:      user.ID,
		})

	})
	JustBeforeEach(func() {
		UpdateTaskBody, _ = json.Marshal(UpdateTaskParams)
	})

	Context("Update Task", func() {
		It("should update task successfully", func() {
			req, err := http.NewRequest("PUT", "/task/"+taskId, bytes.NewBuffer(UpdateTaskBody))
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

			var response types.UpdateTaskResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response.Success).To(BeTrue())
			Expect(response.Message).To(Equal("Task updated successfully"))

			// Assert the task is updated in the database
			task, _ := services.GetTaskRepository(test_helpers.GetTestDatabase()).GetTask(context.TODO(), taskId, user.ID)
			Expect(task.Title).To(Equal("Updated Task"))
			Expect(task.Description).To(Equal("Updated Description"))
			Expect(task.DueDate).To(Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)))
			Expect(int(task.Status)).To(Equal(2))
		})

		When("task does not exist", func() {
			BeforeEach(func() {
				taskId = helpers.GenerateUUID()
			})
			It("should return error", func() {
				req, err := http.NewRequest("PUT", "/task/"+taskId, bytes.NewBuffer(UpdateTaskBody))
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
				req, err := http.NewRequest("PUT", "/task/"+taskId, bytes.NewBuffer(UpdateTaskBody))
				req.Header.Set("Content-Type", "application/json")

				res, err := app.Test(req)

				Expect(err).To(BeNil())
				Expect(res.StatusCode).To(Equal(fiber.StatusNotFound))
			})
		})

	})

})
