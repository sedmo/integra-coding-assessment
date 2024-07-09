package handlers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sedmo/integra-coding-assessment/go-backend/db"
	"github.com/sedmo/integra-coding-assessment/go-backend/db/connectors"
	"github.com/sedmo/integra-coding-assessment/go-backend/handlers"
	"github.com/sedmo/integra-coding-assessment/go-backend/models"
)

func TestHandlers(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Handlers Suite")
}

var _ = ginkgo.Describe("User Handlers", func() {
	var (
		mockDB        *sql.DB
		e             *echo.Echo
		rec           *httptest.ResponseRecorder
		c             echo.Context
		request       *http.Request
		sqlMock       sqlmock.Sqlmock
		mockConnector *connectors.MockConnector
	)

	ginkgo.BeforeEach(func() {
		var err error
		// Initialize mock database
		mockDB, sqlMock, err = sqlmock.New()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		// Initialize Echo and Recorder
		e = echo.New()
		rec = httptest.NewRecorder()

		// Set the DATABASE_URL environment variable before each test
		err = os.Setenv("DATABASE_URL", "postgres://user:password@localhost:5432/mydb?sslmode=disable")
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		// Initialize MockConnector
		mockConnector = &connectors.MockConnector{}
		mockConnector.DB = mockDB
		db.DB = mockDB

		// Initialize database
		db.InitDB(mockConnector)
	})

	ginkgo.AfterEach(func() {
		sqlMock.ExpectClose()
		mockDB.Close()
		os.Unsetenv("DATABASE_URL")
	})

	ginkgo.Describe("CreateUser", func() {
		ginkgo.It("should create a new user", func() {
			user := models.User{
				UserName:   "testuser",
				FirstName:  "Test",
				LastName:   "User",
				Email:      "testuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			sqlMock.ExpectExec("INSERT INTO users").WithArgs(
				user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department,
			).WillReturnResult(sqlmock.NewResult(1, 1))

			err := handlers.CreateUser(c)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusCreated))

			var createdUser models.User
			json.Unmarshal(rec.Body.Bytes(), &createdUser)
			gomega.Expect(createdUser.UserID).NotTo(gomega.BeZero())
			gomega.Expect(createdUser.UserName).To(gomega.Equal(user.UserName))
		})

		ginkgo.It("should return conflict if username already exists", func() {
			user := models.User{
				UserName:   "testuser",
				FirstName:  "Test",
				LastName:   "User",
				Email:      "testuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			sqlMock.ExpectExec("INSERT INTO users").WithArgs(
				user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department,
			).WillReturnError(sqlmock.ErrCancelled) // Simulate conflict error

			err := handlers.CreateUser(c)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusConflict))
		})
	})

	ginkgo.Describe("GetUsers", func() {
		ginkgo.It("should return all users", func() {
			sqlMock.ExpectQuery("SELECT \\* FROM users").WillReturnRows(
				sqlmock.NewRows([]string{"user_id", "user_name", "first_name", "last_name", "email", "user_status", "department"}),
			)

			request = httptest.NewRequest(http.MethodGet, "/users", nil)
			c = e.NewContext(request, rec)

			err := handlers.GetUsers(c)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			var users []models.User
			json.Unmarshal(rec.Body.Bytes(), &users)
			gomega.Expect(users).To(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("UpdateUser", func() {
		ginkgo.It("should update an existing user", func() {
			user := models.User{
				UserName:   "testuser",
				FirstName:  "Test",
				LastName:   "User",
				Email:      "testuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			sqlMock.ExpectExec("INSERT INTO users").WithArgs(
				user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department,
			).WillReturnResult(sqlmock.NewResult(1, 1))

			handlers.CreateUser(c)
			var createdUser models.User
			json.Unmarshal(rec.Body.Bytes(), &createdUser)

			updatedUser := createdUser
			updatedUser.FirstName = "Updated"
			body, _ = json.Marshal(updatedUser)
			request = httptest.NewRequest(http.MethodPut, "/users/"+strconv.FormatInt(createdUser.UserID, 10), bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatInt(createdUser.UserID, 10))

			sqlMock.ExpectExec("UPDATE users").WithArgs(
				updatedUser.UserName, updatedUser.FirstName, updatedUser.LastName, updatedUser.Email, updatedUser.UserStatus, updatedUser.Department, createdUser.UserID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

			err := handlers.UpdateUser(c)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))

			var userAfterUpdate models.User
			json.Unmarshal(rec.Body.Bytes(), &userAfterUpdate)
			gomega.Expect(userAfterUpdate.FirstName).To(gomega.Equal("Updated"))
		})
	})

	ginkgo.Describe("DeleteUser", func() {
		ginkgo.It("should delete a user", func() {
			user := models.User{
				UserName:   "testuser",
				FirstName:  "Test",
				LastName:   "User",
				Email:      "testuser@example.com",
				UserStatus: "A",
				Department: "Engineering",
			}
			body, _ := json.Marshal(user)
			request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c = e.NewContext(request, rec)

			sqlMock.ExpectExec("INSERT INTO users").WithArgs(
				user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department,
			).WillReturnResult(sqlmock.NewResult(1, 1))

			handlers.CreateUser(c)
			var createdUser models.User
			json.Unmarshal(rec.Body.Bytes(), &createdUser)

			request = httptest.NewRequest(http.MethodDelete, "/users/"+strconv.FormatInt(createdUser.UserID, 10), nil)
			c = e.NewContext(request, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.FormatInt(createdUser.UserID, 10))

			sqlMock.ExpectExec("DELETE FROM users").WithArgs(createdUser.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

			err := handlers.DeleteUser(c)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(rec.Code).To(gomega.Equal(http.StatusOK))
		})
	})
})
