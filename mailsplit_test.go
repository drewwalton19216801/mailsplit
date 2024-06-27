package mailsplit

import (
	"os"
	"strings"
	"testing"
)

func TestProcessEmail(t *testing.T) {
	// Sample email content with a base64-encoded attachment
	emailContent := `From: sender@example.com
To: recipient@example.com
Subject: Test email
Content-Type: multipart/mixed; boundary="boundary"

--boundary
Content-Type: text/plain

This is the body of the email.

--boundary
Content-Disposition: attachment; filename="test.txt"
Content-Transfer-Encoding: base64

dGVzdCBjb250ZW50Cg==

--boundary--`

	// Create a temporary directory for attachments
	outputDir, err := os.MkdirTemp("", "attachments")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	// Process the email to save attachments and get the modified email content
	modifiedEmailContent, err := ProcessEmail(emailContent, outputDir)
	if err != nil {
		t.Fatalf("Failed to process email: %v", err)
	}

	// Check that the attachment was saved
	attachmentPath := outputDir + "/test.txt"
	if _, err := os.Stat(attachmentPath); os.IsNotExist(err) {
		t.Fatalf("Attachment was not saved: %v", err)
	}

	// Check that the attachment was removed from the email content
	if strings.Contains(modifiedEmailContent, "Content-Disposition: attachment; filename=\"test.txt\"") {
		t.Fatalf("Attachment was not removed from the email content")
	}

	// Check that the rest of the email content is intact
	if !strings.Contains(modifiedEmailContent, "This is the body of the email.") {
		t.Fatalf("Email content was modified")
	}

	// Clean up the temporary directory
	os.RemoveAll(outputDir)
}
