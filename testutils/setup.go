package testutils

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ouiasy/golang-auth/conf"
	"github.com/ouiasy/golang-auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MockDB creates an in-memory SQLite database for testing
func SetupTestDB(t *testing.T) *gorm.DB {
	// Using PostgreSQL for testing - you may want to use SQLite for unit tests
	// For now, this uses a test database connection
	dsn := "host=localhost user=test password=test dbname=test_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("Skipping test: database connection not available")
	}

	// Auto-migrate the User model
	db.AutoMigrate(&models.User{})

	return db
}

func CreateTestConfig() *conf.GlobalConfiguration {
	godotenv.Load("../.env/.env")
	return &conf.GlobalConfiguration{
		Mail: conf.MailConfiguration{
			ResendApiKey:    os.Getenv("RESEND_API_KEY"),
			ResendFromEmail: os.Getenv("RESEND_EMAIL_FROM"),
		},
		App: conf.AppConfiguration{
			Host: "http://localhost:3000",
		},
	}
}
