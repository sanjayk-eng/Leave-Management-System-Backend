package utils

import (
	"fmt"
	"os"
	"testing"
)

// TestSMTPConfig tests if SMTP configuration can be loaded
func TestSMTPConfig(t *testing.T) {
	// Set test environment variables
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@gmail.com")
	os.Setenv("SMTP_PASSWORD", "testpassword")
	os.Setenv("SMTP_FROM", "test@gmail.com")

	config, err := GetSMTPConfig()
	if err != nil {
		t.Fatalf("Failed to get SMTP config: %v", err)
	}

	if config.Host != "smtp.gmail.com" {
		t.Errorf("Expected host smtp.gmail.com, got %s", config.Host)
	}

	if config.Port != 587 {
		t.Errorf("Expected port 587, got %d", config.Port)
	}

	if config.Username != "test@gmail.com" {
		t.Errorf("Expected username test@gmail.com, got %s", config.Username)
	}

	// Clean up
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")
	os.Unsetenv("SMTP_FROM")
}

// TestSMTPConfigMissing tests error handling when SMTP config is missing
func TestSMTPConfigMissing(t *testing.T) {
	// Ensure no SMTP env vars are set
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")
	os.Unsetenv("SMTP_FROM")

	_, err := GetSMTPConfig()
	if err == nil {
		t.Error("Expected error when SMTP config is missing, but got nil")
	}
}

// TestEmailValidation tests basic email validation
func TestEmailValidation(t *testing.T) {
	// This is a basic test - in a real scenario, you'd want to test with a mock SMTP server
	// or use a test email service
	
	// Set up test config
	os.Setenv("SMTP_HOST", "smtp.gmail.com")
	os.Setenv("SMTP_PORT", "587")
	os.Setenv("SMTP_USERNAME", "test@gmail.com")
	os.Setenv("SMTP_PASSWORD", "testpassword")
	os.Setenv("SMTP_FROM", "test@gmail.com")

	// Test that the function doesn't panic with valid inputs
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SendEmail panicked: %v", r)
		}
	}()

	// This will fail to connect, but should not panic
	err := SendEmail("test@example.com", "Test Subject", "Test Body")
	
	// We expect an error since we're using fake credentials
	if err == nil {
		fmt.Println("Note: Email might have been sent successfully (if real credentials were used)")
	} else {
		fmt.Printf("Expected error occurred: %v\n", err)
	}

	// Clean up
	os.Unsetenv("SMTP_HOST")
	os.Unsetenv("SMTP_PORT")
	os.Unsetenv("SMTP_USERNAME")
	os.Unsetenv("SMTP_PASSWORD")
	os.Unsetenv("SMTP_FROM")
}