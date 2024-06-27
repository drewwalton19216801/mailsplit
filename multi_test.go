package mailsplit

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMultipleAttachments(t *testing.T) {
	emailContent := `From: sender@example.com
To: recipient@example.com
Subject: Test email with multiple attachments
Content-Type: multipart/mixed; boundary="boundary"

--boundary
Content-Type: text/plain

This is the body of the email.

--boundary
Content-Disposition: attachment; filename="file1.txt"
Content-Transfer-Encoding: base64

VGhpcyBpcyBhIHRleHQgYXR0YWNobWVudAo=

--boundary
Content-Disposition: attachment; filename="file2.jpg"
Content-Transfer-Encoding: base64

VGhpcyBpcyBhIGppcGVnIGF0dGFjaG1lbnQKCg==

--boundary
Content-Disposition: attachment; filename="file3.pdf"
Content-Transfer-Encoding: base64

VGhpcyBpcyBhIHBkZiBhdHRhY2htZW50Cg==

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

	// Check that the attachments were saved to the output directory
	expectedAttachments := []string{
		outputDir + "/file1.txt",
		outputDir + "/file2.jpg",
		outputDir + "/file3.pdf",
	}

	for _, expectedAttachment := range expectedAttachments {
		if _, err := os.Stat(expectedAttachment); os.IsNotExist(err) {
			t.Errorf("Expected attachment not found at %s", expectedAttachment)
		}
	}

	// Check that the modified email content does not contain the attachment parts
	attachments := []string{"file1.txt", "file2.jpg", "file3.pdf"}
	for _, attachment := range attachments {
		if strings.Contains(modifiedEmailContent, fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"", attachment)) {
			t.Fatalf("Attachment %s was not removed from the email content", attachment)
		}
	}
}
