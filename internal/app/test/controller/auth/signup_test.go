package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/controller"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/models"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/services"
	test_helpers "github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/test"
	"github.com/Hari-Krishna-Moorthy/task-management-system/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Signup", func() {

	var (
		app      *fiber.App
		userData *types.SignUpRequest
		userJSON []byte
	)

	BeforeEach(func() {
		app = fiber.New()
		app.Post("/signup", controller.GetAuthController(services.GetAuthService(context.TODO(), test_helpers.GetTestDatabase())).SignUp)

		userData = &types.SignUpRequest{
			Username: "hkm",
			Email:    "hkm@email.com",
			Password: "password",
		}

	})

	JustBeforeEach(func() {
		var err error
		userJSON, err = json.Marshal(userData)
		fmt.Printf("userJSON: %v, err: %v\n", userJSON, err)
	})
	AfterEach(func() {
		test_helpers.GetTestClient().Disconnect(context.Background())
	})

	When("request is valid", func() {
		Context("user does not exist", func() {
			It("should create a new user", func() {

				httpReq, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
				httpReq.Header.Set("Content-Type", "application/json")

				res, err := app.Test(httpReq)

				fmt.Printf("res: %v, err: %v\n", res, err)

				// Assert the result
				Expect(err).To(BeNil())
				// Expect(res.StatusCode).To(Equal(fiber.StatusOK))

				var response types.SignUpResponse
				err = json.NewDecoder(res.Body).Decode(&response)
				Expect(err).To(BeNil())
				Expect(response.Success).To(BeTrue())
				Expect(response.Message).To(Equal("User registered successfully"))

				//  Assert the user is created in the database

				user := &models.User{}
				userCollection := test_helpers.GetTestDatabase().Collection("users")
				userCollection.FindOne(context.TODO(), bson.M{"username": userData.Username}).Decode(&user)

				Expect(user.Username).To(Equal(userData.Username))
				Expect(user.Email).To(Equal(userData.Email))
				Expect(bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password))).To(BeNil())

				Expect(user.CreatedAt).NotTo(BeNil())
				Expect(user.UpdatedAt).NotTo(BeNil())
				Expect(user.ID).NotTo(BeNil())

			})
		})
	})

	When("request is invalid", func() {
		Context("user already exists", func() {
			BeforeEach(func() {

				test_helpers.GetTestDatabase().Collection("users").InsertOne(context.TODO(), bson.M{
					"username": "user1",
					"email":    "user1@email.com",
					"password": "user1password",
				})

				userData.Email = "user1@email.com"
				userData.Username = "user1"
			})
			It("should return an error", func() {
				httpReq, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
				httpReq.Header.Set("Content-Type", "application/json")

				res, err := app.Test(httpReq)

				fmt.Printf("res: %v, err: %v\n", res, err)

				// Assert the result
				Expect(err).To(BeNil())
				// Expect(res.StatusCode).To(Equal(fiber.StatusOK))

				var response types.SignUpResponse
				err = json.NewDecoder(res.Body).Decode(&response)
				Expect(err).To(BeNil())
				Expect(response.Success).To(BeFalse())
				Expect(response.Message).To(Equal("user already exists"))
			})
		})

		Context("user details are invalid", func() {
			type TestCase struct {
				Input            types.SignUpRequest
				ExpectedResponse types.RequestError
			}

			var testCases []TestCase

			BeforeEach(func() {
				testCases = []TestCase{
					// Missing Username
					{
						Input: types.SignUpRequest{
							Username: "",
							Email:    "validemail@example.com",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Field: "Username", Tag: "required"},
					},
					// Username Too Short
					{
						Input: types.SignUpRequest{
							Username: "ab",
							Email:    "validemail@example.com",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Field: "Username", Tag: "min", Value: "3"},
					},
					// Username Too Long
					{
						Input: types.SignUpRequest{
							Username: "ThisIsAUsernameThatIsWayTooLongToBeValid",
							Email:    "validemail@example.com",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Field: "Username", Tag: "max", Value: "30"},
					},
					// Missing Email
					{
						Input: types.SignUpRequest{
							Username: "validUsername",
							Email:    "",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Field: "Email", Tag: "required"},
					},
					// Invalid Email Format
					{
						Input: types.SignUpRequest{
							Username: "validUsername",
							Email:    "invalid-email-format",
							Password: "validPassword123",
						},
						ExpectedResponse: types.RequestError{Field: "Email", Tag: "email"},
					},
					// Missing Password
					{
						Input: types.SignUpRequest{
							Username: "validUsername",
							Email:    "validemail@example.com",
							Password: "",
						},
						ExpectedResponse: types.RequestError{Field: "Password", Tag: "required"},
					},
					// Password Too Short
					{
						Input: types.SignUpRequest{
							Username: "validUsername",
							Email:    "validemail@example.com",
							Password: "short",
						},
						ExpectedResponse: types.RequestError{Field: "Password", Tag: "min", Value: "8"},
					},
					// Password Too Long
					{
						Input: types.SignUpRequest{
							Username: "validUsername",
							Email:    "validemail@example.com",
							Password: "ThisPasswordIsWayTooLongAndExceedsTheAllowedLengthLimitOf128Characters" +
								"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
						},
						ExpectedResponse: types.RequestError{Field: "Password", Tag: "max", Value: "128"},
					},
				}

			})
			It("should return an error", func() {

				for _, testCase := range testCases {

					userJSON, _ = json.Marshal(testCase.Input)

					httpReq, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
					httpReq.Header.Set("Content-Type", "application/json")

					res, err := app.Test(httpReq)

					fmt.Printf("res: %v, err: %v\n", res, err)

					// Assert the result
					Expect(err).To(BeNil())
					Expect(res.StatusCode).To(Equal(fiber.StatusNotFound))

					var response []types.RequestError
					err = json.NewDecoder(res.Body).Decode(&response)
					Expect(err).To(BeNil())
					Expect(response[0].Value).To(Equal(testCase.ExpectedResponse.Value))
					Expect(response[0].Field).To(Equal(testCase.ExpectedResponse.Field))
					Expect(response[0].Tag).To(Equal(testCase.ExpectedResponse.Tag))
				}
			})
		})
	})
})
