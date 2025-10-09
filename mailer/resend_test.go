package mailer

import (
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ouiasy/golang-auth/models"
	"github.com/ouiasy/golang-auth/testutils"
)

func createTestUser(t *testing.T) *models.User {
	user := &models.User{
		ID:       uuid.New(),
		Role:     "user",
		Email:    os.Getenv("EMAIL_SEND_TO"), // TODO: change here
		Username: "testuser",
	}
	return user
}

func TestSendConfirmationEmail_Success(t *testing.T) {
	// Skip this test if no API key is available
	// t.Skip("Skipping integration test - requires valid Resend API key")

	db := testutils.SetupTestDB(t)
	config := testutils.CreateTestConfig()
	client := NewEmailClient(config)
	user := createTestUser(t)

	// Create user in database
	db.Create(user)
	defer db.Delete(user)

	maxFreq := 1 * time.Minute
	err := client.SendConfirmationEmail(db, user, user.Email, maxFreq)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify user was updated
	var updatedUser models.User
	db.First(&updatedUser, user.ID)

	if updatedUser.ConfirmationToken == "" {
		t.Error("Expected ConfirmationToken to be set")
	}

	if updatedUser.ConfirmationSentAt == nil {
		t.Error("Expected ConfirmationSentAt to be set")
	}
}

// func TestSendConfirmationEmail_FrequencyLimit(t *testing.T) {
// 	t.Skip("Skipping integration test - requires database connection")

// 	db := setupTestDB(t)
// 	config := createTestConfig()
// 	client := NewEmailClient(config)
// 	user := createTestUser(t)

// 	// Set ConfirmationSentAt to recent time
// 	now := time.Now()
// 	user.ConfirmationSentAt = &now

// 	// Create user in database
// 	db.Create(user)
// 	defer db.Delete(user)

// 	maxFreq := 5 * time.Minute
// 	err := client.SendConfirmationEmail(db, user, user.Email, maxFreq)

// 	if err != ErrorMaxFrequencyLimit {
// 		t.Errorf("Expected ErrorMaxFrequencyLimit, got: %v", err)
// 	}
// }

// func TestSendConfirmationEmail_FrequencyLimitExpired(t *testing.T) {
// 	t.Skip("Skipping integration test - requires database and valid API key")

// 	db := setupTestDB(t)
// 	config := createTestConfig()
// 	client := NewEmailClient(config)
// 	user := createTestUser(t)

// 	// Set ConfirmationSentAt to old time (beyond frequency limit)
// 	oldTime := time.Now().Add(-10 * time.Minute)
// 	user.ConfirmationSentAt = &oldTime

// 	// Create user in database
// 	db.Create(user)
// 	defer db.Delete(user)

// 	maxFreq := 5 * time.Minute
// 	err := client.SendConfirmationEmail(db, user, user.Email, maxFreq)

// 	if err != nil {
// 		t.Errorf("Expected no error when frequency limit expired, got: %v", err)
// 	}
// }

// func TestNewEmailClient(t *testing.T) {
// 	config := createTestConfig()
// 	client := NewEmailClient(config)

// 	if client == nil {
// 		t.Error("Expected client to be created")
// 	}

// 	if client.config != config {
// 		t.Error("Expected config to be set correctly")
// 	}

// 	if client.Client == nil {
// 		t.Error("Expected Resend client to be initialized")
// 	}
// }

// func TestConfirmationEmailData(t *testing.T) {
// 	data := ConfirmationEmailData{
// 		userName:  "testuser",
// 		verifyUrl: "http://localhost:3000/verify?token=abc123",
// 	}

// 	if data.userName != "testuser" {
// 		t.Errorf("Expected userName to be 'testuser', got: %s", data.userName)
// 	}

// 	if data.verifyUrl != "http://localhost:3000/verify?token=abc123" {
// 		t.Errorf("Expected verifyUrl to match, got: %s", data.verifyUrl)
// 	}
// }
