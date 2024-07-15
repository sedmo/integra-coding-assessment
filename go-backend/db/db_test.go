package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/joho/godotenv"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sedmo/integra-coding-assessment/go-backend/db"
	"github.com/sedmo/integra-coding-assessment/go-backend/db/connectors"
)

func TestDB(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "DB Suite")
}

var _ = ginkgo.Describe("DB", func() {
	var (
		mockDB        *sql.DB
		sqlMock       sqlmock.Sqlmock
		mockConnector *connectors.MockConnector
	)

	ginkgo.BeforeEach(func() {
		var err error
		mockConnector = &connectors.MockConnector{}
		mockDB, sqlMock, err = sqlmock.New(sqlmock.MonitorPingsOption(true))
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		mockConnector.DB = mockDB
		mockConnector.Sqlmock = sqlMock

		// Print current working directory
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
		}
		log.Println("Current working directory:", cwd)

		// Adjust the path to the .env file based on the current working directory
		envPath := "../../.env"
		if cwd == "/app" { // Assuming your Docker container sets the working directory to /app
			envPath = ".env"
		}

		// Load environment variables from .env file
		err = godotenv.Load(envPath)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred(), "Failed to load .env file")

		// Verify environment variable is set
		databaseURL := os.Getenv("DATABASE_URL")
		gomega.Expect(databaseURL).ShouldNot(gomega.BeEmpty(), "DATABASE_URL environment variable is not set")
	})

	ginkgo.AfterEach(func() {
		sqlMock.ExpectClose()
		mockDB.Close()
		os.Unsetenv("DATABASE_URL")
	})

	ginkgo.Describe("InitDB", func() {
		ginkgo.It("should successfully connect to the database", func() {
			sqlMock.ExpectPing().WillReturnError(nil)
			db.InitDB(mockConnector)
			gomega.Expect(db.DB.Ping()).ShouldNot(gomega.HaveOccurred())
		})

		ginkgo.It("should retry connection on failure", func() {
			sqlMock.ExpectPing().WillReturnError(sql.ErrConnDone)
			sqlMock.ExpectPing().WillReturnError(nil)
			db.InitDB(mockConnector)
			gomega.Expect(db.DB.Ping()).ShouldNot(gomega.HaveOccurred())
		})
	})

	ginkgo.Describe("RunMigrations", func() {
		ginkgo.It("should run migrations up successfully", func() {
			// Mock migrations up
			sqlMock.ExpectQuery("SELECT CURRENT_DATABASE()").WillReturnRows(sqlmock.NewRows([]string{"database"}).AddRow("testdb"))
			gomega.Expect(mockConnector.RunMigrationsUp(mockDB)).ShouldNot(gomega.HaveOccurred())
		})

		ginkgo.It("should run migrations down successfully", func() {
			// Mock migrations down
			sqlMock.ExpectQuery("SELECT CURRENT_DATABASE()").WillReturnRows(sqlmock.NewRows([]string{"database"}).AddRow("testdb"))
			gomega.Expect(mockConnector.RunMigrationsDown(mockDB)).ShouldNot(gomega.HaveOccurred())
		})
	})
})
