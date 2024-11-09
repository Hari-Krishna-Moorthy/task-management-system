package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Hari-Krishna-Moorthy/task-management-system/config"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/controller"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	test_helpers "github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/test"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/helpers"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("Signin", func() {

	var (
		app  *fiber.App
		user *models.User
	)

	BeforeEach(func() {
		app = fiber.New()
		app.Post("/Signin", controller.GetAuthController(services.GetAuthService(context.TODO(), test_helpers.GetTestDatabase())).SignIn)
	})

	AfterEach(func() {
		test_helpers.GetTestClient().Disconnect(context.Background())
	})

	When("Signin", func() {
		Context("user details are invalid", func() {
			type TestCase struct {
				Input            types.SignInRequest
				ExpectedResponse types.RequestError
				Success          bool
				Validate         bool
			}

			var testCases []TestCase

			BeforeEach(func() {

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

				testCases = []TestCase{
					// Valid Input with Email
					{
						Input: types.SignInRequest{
							Email:    "validemail@example.com",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Value: "user signed in successfully"},
						Success:          true,
						Validate:         true,
					},
					// Valid Input with Username
					{
						Input: types.SignInRequest{
							Username: "validUsername",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Value: "user signed in successfully"},
						Validate:         true,
						Success:          true,
					},
					// Both Email and Username Provided (valid)
					{
						Input: types.SignInRequest{
							Email:    "validemail@example.com",
							Username: "validUsername",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Value: "user signed in successfully"},
						Validate:         true,
						Success:          true,
					},
					// Missing Both Email and Username
					{
						Input: types.SignInRequest{
							Email:    "",
							Username: "",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Tag: "required_without", Field: "Email", Value: "Username"},
					},
					// Invalid Email Format with No Username
					{
						Input: types.SignInRequest{
							Email:    "invalid-email",
							Username: "",
							Password: "validPassword123",
						},
						Validate:         true,
						Success:          false,
						ExpectedResponse: types.RequestError{Value: "user not found"},
					},
					// Username Too Short, No Email
					{
						Input: types.SignInRequest{
							Email:    "",
							Username: "ab",
							Password: "validPassword123",
						},
						Validate:         true,
						Success:          false,
						ExpectedResponse: types.RequestError{Value: "user not found"},
					},
					// Username Too Long, No Email
					{
						Input: types.SignInRequest{
							Email:    "",
							Username: "ThisIsAUsernameThatIsWayTooLongToBeValid",
							Password: "validPassword123",
						},
						Validate:         true,
						Success:          false,
						ExpectedResponse: types.RequestError{Value: "user not found"},
					},
					// Missing Password
					{
						Input: types.SignInRequest{
							Email:    "validemail@example.com",
							Username: "",
							Password: "",
						},
						ExpectedResponse: types.RequestError{Field: "Password", Tag: "required"},
					},
					// Password Too Short
					{
						Input: types.SignInRequest{
							Email:    "validemail@example.com",
							Username: "",
							Password: "short",
						},
						ExpectedResponse: types.RequestError{Value: "8", Field: "Password", Tag: "min"},
					},
					// Password Too Long
					{
						Input: types.SignInRequest{
							Email:    "validemail@example.com",
							Username: "",
							Password: "ThisPasswordIsWayTooLongAndExceedsTheAllowedLengthLimitOf128Characters" +
								"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
						},
						ExpectedResponse: types.RequestError{Value: "128", Field: "Password", Tag: "max"},
					},
				}

			})

			It("should return an error response", func() {

				for _, testCase := range testCases {

					log.Printf("\n Running Test Case input: %+v, expected response: %+v\n", testCase.Input, testCase.ExpectedResponse)

					userJSON, _ := json.Marshal(testCase.Input)

					httpReq, _ := http.NewRequest("POST", "/signin", bytes.NewBuffer(userJSON))
					httpReq.Header.Set("Content-Type", "application/json")

					res, err := app.Test(httpReq)

					if testCase.Validate {

						Expect(err).To(BeNil())
						// Expect(res.StatusCode).To(Equal(fiber.StatusOK))

						var response types.SignInResponse
						err = json.NewDecoder(res.Body).Decode(&response)
						Expect(err).To(BeNil())
						// Expect(response.Success).To(Equal(testCase.Success))
						Expect(response.Message).To(Equal(testCase.ExpectedResponse.Value))

						if testCase.Success {

							cookie := res.Header.Get("Set-Cookie")
							cookieParts := strings.Split(cookie, ";")
							Expect(cookieParts[0]).To(HavePrefix("token="))
							tokenString := strings.Split(cookieParts[0], "=")[1]

							token, err := jwt.ParseWithClaims(tokenString, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
								return []byte(config.GetConfig().Auth.JWTSecret), nil
							})
							Expect(err).To(BeNil())

							// Extract the claims
							claims, ok := token.Claims.(*types.JWTClaims)
							log.Println("claims: ", claims)
							Expect(ok).To(BeTrue())

							// Validate user ID
							Expect(user.ID).To((ContainSubstring(claims.UserID)))
							Expect(claims.Username).To(Equal("validUsername"))
							Expect(claims.Email).To(Equal("validemail@example.com"))

							// Validate expiration date: token should expire within 3 days
							expireAt := time.Unix(claims.ExpireAt, 0)
							Expect(expireAt.Before(time.Now().Add(3 * 24 * time.Hour))).To(BeTrue())
						}

					} else {
						Expect(err).To(BeNil())
						Expect(res.StatusCode).To(Equal(fiber.StatusNotFound))
						var response []types.RequestError
						err = json.NewDecoder(res.Body).Decode(&response)
						Expect(err).To(BeNil())
						Expect(response[0].Value).To(Equal(testCase.ExpectedResponse.Value))
						Expect(response[0].Field).To(Equal(testCase.ExpectedResponse.Field))
						Expect(response[0].Tag).To(Equal(testCase.ExpectedResponse.Tag))
					}
				}
			})
		})
	})

})
