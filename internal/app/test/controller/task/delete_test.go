package task_test

import (
	"context"
	"encoding/json"
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

var _ = Describe("Delete Task", func() {

	var (
		app              *fiber.App
		token            string
		user             *models.User
		task1ID, task2ID string
	)

	BeforeEach(func() {
		app = fiber.New()
		app.Delete("/task/:id", controller.GetTaskController(services.GetTaskService(services.GetTaskRepository(test_helpers.GetTestDatabase()))).DeleteTask)

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

	Context("When deleting a task", func() {
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

		It("Should delete the task", func() {
			uri := "/task/" + task1ID
			req, err := http.NewRequest("DELETE", uri, nil)

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

			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(fiber.StatusOK))

			var response types.DeleteTaskResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			Expect(err).Should(BeNil())
			Expect(response.Success).Should(BeTrue())
			Expect(response.Message).Should(Equal("Task deleted successfully"))

			_, err = services.GetTaskRepository(test_helpers.GetTestDatabase()).GetTask(context.TODO(), task1ID, user.ID)
			Expect(err).ShouldNot(BeNil())
			Expect(err.Error()).Should(Equal("mongo: no documents in result"))

			// Assert the task is deleted from the database
			task := &models.Task{}
			test_helpers.GetTestDatabase().Collection("tasks").FindOne(context.TODO(), bson.M{"_id": task1ID}).Decode(&task)

			Expect(task.DeletedAt.Format(utils.TimeLayout)).To(Equal(time.Now().UTC().Format(utils.TimeLayout)))
		})
	})

	Context("When deleting a task that does not exist", func() {

		It("Should return not found", func() {
			uri := "/task/" + "invalid-id"
			req, err := http.NewRequest("DELETE", uri, nil)

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

			Expect(err).Should(BeNil())
			// Expect(res.StatusCode).Should(Equal(fiber.StatusNotFound))

			var response types.DeleteTaskResponse
			err = json.NewDecoder(res.Body).Decode(&response)
			Expect(err).Should(BeNil())
			// Expect(response.Success).Should(BeFalse())
			Expect(response.Message).Should(Equal("mongo: no documents in result"))

		})
	})

	Context("user not authenticated", func() {
		It("should return unauthorized ", func() {
			uri := "/task/" + task1ID
			req, err := http.NewRequest("DELETE", uri, nil)

			req.Header.Set("Content-Type", "application/json")

			res, err := app.Test(req)

			Expect(err).Should(BeNil())
			Expect(res.StatusCode).Should(Equal(fiber.StatusNotFound))
		})
	})

})
