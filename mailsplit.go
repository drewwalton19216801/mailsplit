package mailsplit

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"os"
	"strings"
)

// saveAttachments parses the email and saves base64 encoded attachments to disk
func saveAttachments(emailContent string, outputDir string) error {
	msg, err := mail.ReadMessage(strings.NewReader(emailContent))
	if err != nil {
		return fmt.Errorf("failed to parse email: %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return fmt.Errorf("failed to parse content type: %v", err)
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return fmt.Errorf("email is not a multipart email")
	}

	mr := multipart.NewReader(msg.Body, params["boundary"])
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read next part: %v", err)
		}

		disposition := p.Header.Get("Content-Disposition")
		if disposition == "" {
			continue
		}

		mediaType, params, err := mime.ParseMediaType(disposition)
		if err != nil {
			log.Printf("Failed to parse media type: %v", err)
			continue
		}

		if !strings.HasPrefix(mediaType, "attachment") {
			continue
		}

		filename := params["filename"]
		if filename == "" {
			continue
		}

		filename = strings.ReplaceAll(filename, " ", "_")
		filename = strings.NewReplacer("<", "", ">", "", ":", "_", "\"", "_", "/", "_", "\\", "_", ",", "_").Replace(filename)

		attachmentData, err := io.ReadAll(p)
		if err != nil {
			return fmt.Errorf("failed to read attachment data: %v", err)
		}

		decodedData := make([]byte, base64.StdEncoding.DecodedLen(len(attachmentData)))
		n, err := base64.StdEncoding.Decode(decodedData, attachmentData)
		if err != nil {
			return fmt.Errorf("failed to decode attachment data: %v", err)
		}

		decodedData = decodedData[:n]

		outputPath := fmt.Sprintf("%s/%s", outputDir, filename)

		err = os.WriteFile(outputPath, decodedData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write attachment to disk: %v", err)
		}

		log.Printf("Saved attachment: %s", outputPath)
	}

	return nil
}

// removeAttachments reconstructs the email without attachments
func removeAttachments(emailContent string) (string, error) {
	msg, err := mail.ReadMessage(strings.NewReader(emailContent))
	if err != nil {
		return "", fmt.Errorf("failed to parse email: %v", err)
	}

	mediaType, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("failed to parse media type: %v", err)
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return "", fmt.Errorf("email is not multipart")
	}

	var newEmailContent bytes.Buffer
	// Reconstruct the headers without attachments
	newEmailContent.WriteString(fmt.Sprintf("From: %s\n", msg.Header.Get("From")))
	newEmailContent.WriteString(fmt.Sprintf("To: %s\n", msg.Header.Get("To")))
	newEmailContent.WriteString(fmt.Sprintf("Subject: %s\n", msg.Header.Get("Subject")))
	newEmailContent.WriteString(fmt.Sprintf("Content-Type: %s\n", msg.Header.Get("Content-Type")))
	newEmailContent.WriteString(fmt.Sprintf("MIME-Version: %s\n", msg.Header.Get("MIME-Version")))

	// Write the boundary for the multipart content
	newEmailContent.WriteString(fmt.Sprintf("Content-Type: %s; boundary=\"%s\"\n\n", mediaType, params["boundary"]))

	// Write any other headers that were in the original email
	for k, v := range msg.Header {
		if k != "From" && k != "To" && k != "Subject" && k != "Content-Type" && k != "MIME-Version" {
			newEmailContent.WriteString(fmt.Sprintf("%s: %s\n", k, v[0]))
		}
	}

	// Create a new multipart writer for the reconstructed content
	mw := multipart.NewWriter(&newEmailContent)
	defer mw.Close()

	mr := multipart.NewReader(msg.Body, params["boundary"])

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("failed to get next part: %v", err)
		}

		disposition := p.Header.Get("Content-Disposition")
		if disposition == "" {
			pw, err := mw.CreatePart(p.Header)
			if err != nil {
				return "", fmt.Errorf("failed to create part: %v", err)
			}
			if _, err := io.Copy(pw, p); err != nil {
				return "", fmt.Errorf("failed to copy part: %v", err)
			}
			continue
		}

		mediaType, params, err := mime.ParseMediaType(disposition)
		if err != nil {
			return "", fmt.Errorf("failed to parse media type: %v", err)
		}

		if !strings.HasPrefix(mediaType, "attachment") {
			pw, err := mw.CreatePart(p.Header)
			if err != nil {
				return "", fmt.Errorf("failed to create part: %v", err)
			}
			if _, err := io.Copy(pw, p); err != nil {
				return "", fmt.Errorf("failed to copy part: %v", err)
			}
		}

		if strings.HasPrefix(mediaType, "attachment") {
			continue
		}

		filename := params["filename"]
		if filename == "" {
			pw, err := mw.CreatePart(p.Header)
			if err != nil {
				return "", fmt.Errorf("failed to create part: %v", err)
			}
			if _, err := io.Copy(pw, p); err != nil {
				return "", fmt.Errorf("failed to copy part: %v", err)
			}
		}
	}

	return newEmailContent.String(), nil
}

// ProcessEmail coordinates saving attachments and removing them from the email
func ProcessEmail(emailContent string, outputDir string) (string, error) {
	// Save the attachments
	err := saveAttachments(emailContent, outputDir)
	if err != nil {
		return "", err
	}

	newEmailContent, err := removeAttachments(emailContent)
	if err != nil {
		return "", err
	}

	return newEmailContent, nil
}
