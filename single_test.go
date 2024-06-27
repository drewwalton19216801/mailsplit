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
	defer func() {
		if err := os.RemoveAll(outputDir); err != nil {
			t.Logf("Failed to remove temporary directory: %v", err)
		}
	}()

	modifiedEmailContent, err := ProcessEmail(emailContent, outputDir)
	if err != nil {
		t.Fatalf("Failed to process email: %v", err)
	}

	// Check that the attachment was saved to the output directory
	expectedAttachmentPath := outputDir + "/test.txt"
	if _, err := os.Stat(expectedAttachmentPath); os.IsNotExist(err) {
		t.Errorf("Expected attachment not found at %s", expectedAttachmentPath)
	}

	// Split the modified email content into parts
	parts := strings.Split(modifiedEmailContent, "\n\n")

	// Check that the modified email content has the expected number of parts
	if len(parts) < 2 {
		t.Errorf("Modified email content has unexpected number of parts")
	}

	// Check the header part
	headerPart := strings.TrimSpace(parts[0])
	if !strings.HasPrefix(headerPart, "From: sender@example.com") ||
		!strings.Contains(headerPart, "Subject: Test email with single attachment") ||
		!strings.Contains(headerPart, "Content-Type: multipart/mixed; boundary=\"boundary\"") {
		t.Errorf("Modified email header does not match expected content")
		// Print the header part for debugging purposes
		t.Logf("Modified email header: %s", headerPart)
	}

	// Check the body part
	bodyPart := strings.TrimSpace(parts[1])
	if !strings.Contains(bodyPart, "This is the body of the email.") {
		t.Errorf("Modified email body does not match expected content")
	}

	// Check that there are no additional parts beyond the expected ones
	if len(parts) > 3 {
		t.Errorf("Modified email content contains unexpected MIME parts")
	}

	// Check that the attachment part is removed
	if strings.Contains(bodyPart, "Content-Disposition: attachment; filename=\"test.txt\"") {
		t.Errorf("Attachment part was not removed from the email body")
	}
}
