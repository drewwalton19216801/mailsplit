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
Content-Disposition: attachment; filename="test1.txt"
Content-Transfer-Encoding: base64

VGhpcyBpcyBhIHNpbmdsZSBhdHRhY2htZW50Cg==

--boundary
Content-Disposition: attachment; filename="test2.pdf"
Content-Transfer-Encoding: base64

JVBERi0xLjUNCiW1tbW1DQoxIDAgb2JqDQo8PCAvVHlwZSAvUGFnZQ0KL0tpZHMgWzQgMCBSXS9J
RCBvYmoNCjw8IC9UeXBlIC9QYWdlcw0KL01lZGlhQm94IFswIDAgNTk1LjIyOSA3MTIuOTddDQov
TWVkaWFCb3ggWzAgMCA1OTUuMjI5IDcxMi45N10NCi9Db250ZW50cyA1IDAgUg0KL1Jlc291cmNl
cyA8PA0KL0ZpbHRlciAvRmxhdGVEZWNvZGUNCj4NCnN0cmVhbQ0KPDwvU2l6ZSAyL0ZpbHRlciAv
RmxhdGVEZWNvZGUNCi9MZW5ndGggNDYNCi9Gb250IDw8IC9GMSAxMiAvRjIgMiBGMiA3IDAgUiA+
Pj4NCnN0YXJ0eHJlZg0KMjcyNA0KJSVFT0YNCg==

--boundary--
`

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
		outputDir + "/test1.txt",
		outputDir + "/test2.pdf",
	}

	for _, expectedAttachment := range expectedAttachments {
		if _, err := os.Stat(expectedAttachment); os.IsNotExist(err) {
			t.Errorf("Expected attachment not found at %s", expectedAttachment)
		}
	}

	// Check that the modified email content does not contain the attachment parts
	attachments := []string{"test1.txt", "test2.pdf"}
	for _, attachment := range attachments {
		if strings.Contains(modifiedEmailContent, fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"", attachment)) {
			t.Fatalf("Attachment %s was not removed from the email content", attachment)
		}
	}
}
