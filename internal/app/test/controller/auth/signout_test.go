package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/controller"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	test_helpers "github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/test"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/utils"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
)

var _ = Describe("AuthService SignOut", func() {
	var app *fiber.App

	BeforeEach(func() {
		app = fiber.New()
		app.Post("/signout", controller.GetAuthController(services.GetAuthService(context.TODO(), test_helpers.GetTestDatabase())).SignOut)

	})

	Context("when the cookie is present", func() {
		It("should sign out successfully and remove the cookie", func() {
			// Create a request with the token cookie set

			httpReq, err := http.NewRequest("POST", "/signout", nil)
			httpReq.Header.Set("Content-Type", "application/json")

			cookie := &http.Cookie{
				Name:     utils.CookieKeyToken,
				Value:    "dummy_token",
				Expires:  time.Now().Add(utils.JWT_TOKEN_EXPIRY),
				HttpOnly: true,
				Domain:   config.GetConfig().Server.AppDomain,
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
				Secure:   true,
				MaxAge:   int(utils.JWT_TOKEN_EXPIRY.Seconds()),
			}

			httpReq.AddCookie(cookie)

			// Create a response recorder to capture the response
			resp, err := app.Test(httpReq)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var response types.SignOutResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response.Success).To(BeTrue())
			Expect(response.Message).To(Equal("Signed out successfully"))
		})
	})

	Context("when the cookie is absent", func() {
		It("should return a message that the user is not signed in", func() {
			// Create a request without the token cookie
			httpReq, err := http.NewRequest("POST", "/signout", nil)
			httpReq.Header.Set("Content-Type", "application/json")

			// Create a response recorder to capture the response
			resp, err := app.Test(httpReq)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

			var response types.SignOutResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			Expect(err).To(BeNil())
			Expect(response.Success).To(BeFalse())
			Expect(response.Message).To(Equal("user not signed in"))
		})
	})
})
