package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/drewwalton19216801/mailsplit/mailsplit"
)

// mailsplit is a command line tool that processes an email and saves attachments to disk
// Usage: mailsplit <path_to_email_file>

func main() {
	// Read email from a file
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	emailContent, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Create the output directory if it doesn't exist
	outputDir := "attachments"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.Mkdir(outputDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Process the email to save attachments and get the modified email content
	newEmailContent, err := mailsplit.ProcessEmail(string(emailContent), outputDir)
	if err != nil {
		log.Fatalf("Failed to process email: %v", err)
	}

	// Save the modified email content to a new file
	outputFile, err := os.Create("output.eml")
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()
	_, err = outputFile.WriteString(newEmailContent)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Email processed successfully!")
}
