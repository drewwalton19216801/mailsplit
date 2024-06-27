# MailSplit

MailSplit is a Go library for parsing emails, saving base64-encoded attachments to disk, and reconstructing the email without the attachments.

## Installation

To use the MailSplit library, you need to have [Go](https://golang.org/) installed on your machine.

1. Create a new directory for your Go project and navigate into it:

    ```sh
    mkdir mailsplit-project
    cd mailsplit-project
    ```

2. Initialize a new Go module:

    ```sh
    go mod init mailsplit-project
    ```

3. Install the `mailsplit` package in your project:

    ```sh
    go get github.com/drewwalton19216801/mailsplit
    ```

## Usage

To use the MailSplit library, you can create a simple Go program. Here is an example:

```go
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "github.com/drewwalton19216801/mailsplit"
)

func main() {
    // Read the email from a file
    emailFile, err := os.Open("example_email.eml")
    if err != nil {
        log.Fatalf("Failed to open email: %v", err)
    }
    defer emailFile.Close()

    emailContent, err := io.ReadAll(emailFile)
    if err != nil {
        log.Fatalf("Failed to read email: %v", err)
    }

    // Create the output directory if it does not exist
    outputDir := "attachments"
    if _, err := os.Stat(outputDir); os.IsNotExist(err) {
        err := os.Mkdir(outputDir, 0755)
        if err != nil {
            log.Fatalf("Failed to create output directory: %v", err)
        }
    }

    // Process the email to save attachments and get the modified email content
    newEmailContent, err := mailsplit.ProcessEmail(string(emailContent), outputDir)
    if err != nil {
        log.Fatalf("Failed to process email: %v", err)
    }

    // Save the new email content to a file
    err = os.WriteFile("modified_email.eml", []byte(newEmailContent), 0644)
    if err != nil {
        log.Fatalf("Failed to save modified email: %v", err)
    }

    fmt.Println("Attachments saved and email modified successfully.")
}
```

## Building and Running the Example Program

To build and run the included example program, titled `mailsplit`, you can use the following commands:

```sh
go build ./cmd/mailsplit
./mailsplit /path/to/email.eml
```

Any attachments will be saved to the `attachments` directory.

## Contributing

If you'd like to contribute to the project, please feel free to open an issue or pull request.

## License

This project is licensed under the MIT License.