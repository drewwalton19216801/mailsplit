package mailsplit

import (
	"os"
	"strings"
	"testing"
)

func TestSingleAttachment(t *testing.T) {
	emailContent := `From: sender@example.com
To: recipient@example.com
Subject: Test email with single attachment
Content-Type: multipart/mixed; boundary="boundary"

--boundary
Content-Type: text/plain

This is the body of the email.

--boundary
Content-Disposition: attachment; filename="test.txt"
Content-Transfer-Encoding: base64

VGhpcyBpcyBhIHNpbmdsZSBhdHRhY2htZW50Cg==

--boundary--`

	outputDir, err := os.MkdirTemp("", "mailsplit")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(outputDir)

	modifiedEmailContent, err := ProcessEmail(emailContent, outputDir)
	if err != nil {
		t.Fatalf("Failed to process email: %v", err)
	}

	// Check that the attachment was saved to the output directory
	expectedAttachmentPath := outputDir + "/test.txt"
	if _, err := os.Stat(expectedAttachmentPath); os.IsNotExist(err) {
		t.Errorf("Expected attachment not found at %s", expectedAttachmentPath)
	}

	// Normalize line endings for comparison
	modifiedEmailContentNormalized := strings.ReplaceAll(modifiedEmailContent, "\r\n", "\n")

	// Split content by boundary and check the relevant parts
	parts := strings.Split(modifiedEmailContentNormalized, "\n--boundary")
	if len(parts) < 1 { // 1 part = 2 boundary lines
		t.Errorf("Modified email content does not contain expected boundary")
	}

	// Check the body part
	bodyPart := strings.TrimSpace(parts[0])
	if !strings.Contains(bodyPart, "This is the body of the email.") {
		t.Errorf("Modified email body does not match expected content")
	}

	// Check that there are no additional parts beyond the expected ones
	if len(parts) > 1 {
		t.Errorf("Modified email content contains unexpected MIME parts")
	}

	// Check that the attachment part is removed
	if strings.Contains(bodyPart, "Content-Disposition: attachment; filename=\"test.txt\"") {
		t.Errorf("Attachment part was not removed from the email body")
	}
}
