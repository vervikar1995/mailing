package utils

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"regexp"
	"sync"
)

const maxConcurrency = 20 // Maximum number of concurrently running goroutines

func SendEmail(to []string, subject, body string) error {
	// SMTP server configuration
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	sender := "" //your gmail account
	password := "" // App-specific password for Gmail

	// Compose the message
	msg := "From: " + sender + "\n" +
		"To: " + to[0] + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	// Set up authentication
	auth := smtp.PlainAuth("", sender, password, smtpHost)

	// Send the email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, to, []byte(msg))
	if err != nil {
		return err
	}

	fmt.Println("Email sent successfully:", to[0])
	return nil
}

func GetValidEmails() ([]string, error) {
	// Open file containing email addresses
	file, err := os.Open("db.txt")
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Regex for email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	validEmails := []string{}
	scanner := bufio.NewScanner(file)

	// Check each email and add valid ones to the list
	for scanner.Scan() {
		email := scanner.Text()
		if emailRegex.MatchString(email) {
			validEmails = append(validEmails, email)
		} else {
			fmt.Println("Invalid email:", email)
		}
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return validEmails, nil
}

func SendEmailsToValidRecipients(validEmails []string) error {
	// Open the file with the message text
	messageFile, err := os.Open("letter.txt")
	if err != nil {
		return fmt.Errorf("error opening message file: %v", err)
	}
	defer messageFile.Close()

	// Read all lines from the file into the `body` variable
	var body string
	scanner := bufio.NewScanner(messageFile)
	for scanner.Scan() {
		body += scanner.Text() + "\n"
	}

	// Handle any scanning error
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading message file: %v", err)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxConcurrency)

	// Launch goroutines to send emails to each recipient
	for _, email := range validEmails {
		wg.Add(1)
		semaphore <- struct{}{} // Limit the number of concurrently running goroutines

		go func(recipient string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release slot for the next goroutine

			subject := "Wisespace"
			err := SendEmail([]string{recipient}, subject, body)
			if err != nil {
				log.Printf("Error sending email to %s: %v", recipient, err)
			}
		}(email)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return nil
}
