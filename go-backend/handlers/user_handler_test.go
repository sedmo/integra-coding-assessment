package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sedmo/integra-coding-assessment/go-backend/db"
	"github.com/sedmo/integra-coding-assessment/go-backend/db/connectors"
	"github.com/sedmo/integra-coding-assessment/go-backend/handlers"
	"github.com/sedmo/integra-coding-assessment/go-backend/models"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = Describe("User Handlers", func() {
	var (
		e             *echo.Echo
		rec           *httptest.ResponseRecorder
		c             echo.Context
		request       *http.Request
		mockConnector *connectors.MockConnector
	)

	BeforeEach(func() {
		var err error
		// Initialize MockConnector
		mockConnector = &connectors.MockConnector{}
		_, err = mockConnector.Open("")
		Expect(err).ShouldNot(HaveOccurred())

		// Initialize Echo and Recorder
		e = echo.New()
		rec = httptest.NewRecorder()

		// Set the DATABASE_URL environment variable before each test
		err = os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
		Expect(err).NotTo(HaveOccurred())

		// Set up expectations for InitDB
		mockConnector.Sqlmock.ExpectQuery("SELECT CURRENT_DATABASE()").WillReturnRows(sqlmock.NewRows([]string{"current_database"}).AddRow("testdb"))
		mockConnector.Sqlmock.ExpectQuery("SELECT (.+) FROM information_schema.tables WHERE table_schema = 'public'").WillReturnRows(sqlmock.NewRows([]string{}))
		mockConnector.Sqlmock.ExpectQuery("SELECT (.+) FROM information_schema.tables WHERE table_schema = 'public'").WillReturnRows(sqlmock.NewRows([]string{}))

		// Initialize database
		err = db.InitDB(mockConnector)
		Expect(err).NotTo(HaveOccurred())

		// Clear expectations set by InitDB
		mockConnector.Sqlmock.ExpectationsWereMet()
	})

	AfterEach(func() {
		mockConnector.Sqlmock.ExpectClose()
		mockConnector.DB.Close()
		os.Unsetenv("DATABASE_URL")
	})

	Describe("CreateUser", func() {
		It("should return error for invalid user data", func() {
			// Simulate sending invalid user data (empty body)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("{}")))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			// Call CreateUser handler
			err := handlers.CreateUser(c)

			// Expect BadRequest due to validation failure
			Expect(err).Should(BeAssignableToTypeOf(&echo.HTTPError{}))
			httpErr := err.(*echo.HTTPError)
			Expect(httpErr.Code).To(Equal(http.StatusBadRequest))
		})

		It("should check for existing username and return conflict if exists", func() {
			user := models.User{
				UserName:   "existinguser",
				FirstName:  "Existing",
				LastName:   "User",
				Email:      "existinguser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}

			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)
			// Expect a query to check if username exists
			mockConnector.Sqlmock.ExpectQuery("SELECT \\* FROM users WHERE user_name = \\$1").
				WithArgs(user.UserName).
				WillReturnRows(sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "user_status", "department"}).
					AddRow(1, user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department))

			err := handlers.CreateUser(c)

			// Expect StatusConflict due to existing username
			Expect(err).Should(BeAssignableToTypeOf(&echo.HTTPError{}))
			httpErr := err.(*echo.HTTPError)
			Expect(httpErr.Code).To(Equal(http.StatusConflict))
		})

		It("should successfully create a new user", func() {
			user := models.User{
				UserName:   "newuser",
				FirstName:  "New",
				LastName:   "User",
				Email:      "newuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			// Expect a query to check if the username already exists
			mockConnector.Sqlmock.ExpectQuery("SELECT \\* FROM users WHERE user_name = \\$1").
				WithArgs(user.UserName).
				WillReturnRows(sqlmock.NewRows([]string{})) // No existing user

			// Expect the insert query
			mockConnector.Sqlmock.ExpectQuery("INSERT INTO users \\(user_name,first_name,last_name,email,user_status,department\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING user_id").
				WithArgs(user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department).
				WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(1))

			// Call CreateUser handler
			err := handlers.CreateUser(c)

			// Expect no error and a StatusCreated
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusCreated))

			// Verify that the response contains the user ID
			var responseUser models.User
			err = json.Unmarshal(rec.Body.Bytes(), &responseUser)
			Expect(err).ToNot(HaveOccurred())
			Expect(responseUser.UserID).To(Equal(int64(1)))
		})
	})

	Describe("GetUsers", func() {
		It("should return all users", func() {
			// Mock the database response
			mockConnector.Sqlmock.ExpectQuery("SELECT \\* FROM users").WillReturnRows(
				sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "user_status", "department"}).
					AddRow(1, "user1", "User", "One", "user1@example.com", "A", "Engineering").
					AddRow(2, "user2", "User", "Two", "user2@example.com", "I", "Marketing"),
			)

			// Make the request
			request = httptest.NewRequest(http.MethodGet, "/users", nil)
			c = e.NewContext(request, rec)

			// Call the handler
			err := handlers.GetUsers(c)

			// Expect no error
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			// Verify the response
			var users []models.User
			err = json.Unmarshal(rec.Body.Bytes(), &users)
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(2))
			Expect(users[0].UserName).To(Equal("user1"))
			Expect(users[1].UserName).To(Equal("user2"))
		})
	})

	Describe("UpdateUser", func() {
		It("should update an existing user", func() {
			user := models.User{
				UserName:   "newuser",
				FirstName:  "New",
				LastName:   "User",
				Email:      "newuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			// Expect a query to check if the username already exists
			mockConnector.Sqlmock.ExpectQuery("SELECT \\* FROM users WHERE user_name = \\$1").
				WithArgs(user.UserName).
				WillReturnRows(sqlmock.NewRows([]string{})) // No existing user

			// Expect the insert query
			mockConnector.Sqlmock.ExpectQuery("INSERT INTO users \\(user_name,first_name,last_name,email,user_status,department\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING user_id").
				WithArgs(user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department).
				WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(int64(1)))

			// Call CreateUser handler
			err := handlers.CreateUser(c)

			// Expect no error and a StatusCreated
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			var createdUser models.User
			json.Unmarshal(rec.Body.Bytes(), &createdUser)

			// reset rec before updating the user to ensure no previous response is present
			rec = httptest.NewRecorder()

			updatedUser := models.User{
				UserID:     createdUser.UserID,
				UserName:   "newuser",
				FirstName:  "Updated",
				LastName:   "User",
				Email:      "newuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ = json.Marshal(updatedUser)
			request = httptest.NewRequest(http.MethodPut, "/users/"+strconv.FormatInt(createdUser.UserID, 10), bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatInt(createdUser.UserID, 10))

			// Expect a query to check if the username already exists
			mockConnector.Sqlmock.ExpectQuery("SELECT \\* FROM users WHERE user_name = \\$1").
				WithArgs(updatedUser.UserName).
				WillReturnRows(sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "user_status", "department"}).
					AddRow(int64(1), updatedUser.UserName, updatedUser.FirstName, updatedUser.LastName, updatedUser.Email, updatedUser.UserStatus, updatedUser.Department))

			mockConnector.Sqlmock.ExpectExec("UPDATE users SET user_name = \\$1, first_name = \\$2, last_name = \\$3, email = \\$4, user_status = \\$5, department = \\$6 WHERE user_id = \\$7").
				WithArgs(updatedUser.UserName, updatedUser.FirstName, updatedUser.LastName, updatedUser.Email, updatedUser.UserStatus, updatedUser.Department, strconv.FormatInt(createdUser.UserID, 10)).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = handlers.UpdateUser(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))

			var userAfterUpdate models.User
			json.Unmarshal(rec.Body.Bytes(), &userAfterUpdate)
			Expect(userAfterUpdate.FirstName).To(Equal("Updated"))
		})
	})

	Describe("DeleteUser", func() {
		It("should delete a user", func() {
			user := models.User{
				UserName:   "newuser",
				FirstName:  "New",
				LastName:   "User",
				Email:      "newuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			// Expect a query to check if the username already exists
			mockConnector.Sqlmock.ExpectQuery("SELECT \\* FROM users WHERE user_name = \\$1").
				WithArgs(user.UserName).
				WillReturnRows(sqlmock.NewRows([]string{})) // No existing user

			// Expect the insert query
			mockConnector.Sqlmock.ExpectQuery("INSERT INTO users \\(user_name,first_name,last_name,email,user_status,department\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5,\\$6\\) RETURNING user_id").
				WithArgs(user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department).
				WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(int64(1)))

			// Call CreateUser handler
			err := handlers.CreateUser(c)

			// Expect no error and a StatusCreated
			Expect(err).ToNot(HaveOccurred())
			Expect(rec.Code).To(Equal(http.StatusCreated))
			var createdUser models.User
			json.Unmarshal(rec.Body.Bytes(), &createdUser)

			// Ensure all previous expectations were met
			err = mockConnector.Sqlmock.ExpectationsWereMet()
			Expect(err).ShouldNot(HaveOccurred())

			// reset rec before deleting the user to ensure no previous response is present
			rec = httptest.NewRecorder()

			// Set new expectations for the delete operation
			request = httptest.NewRequest(http.MethodDelete, "/users/"+strconv.FormatInt(createdUser.UserID, 10), nil)
			c = e.NewContext(request, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatInt(createdUser.UserID, 10))

			mockConnector.Sqlmock.ExpectExec("DELETE FROM users WHERE user_id = \\$1").
				WithArgs(strconv.FormatInt(createdUser.UserID, 10)).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = handlers.DeleteUser(c)
			Expect(err).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
		})
	})

})
