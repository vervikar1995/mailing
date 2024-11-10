package main

import (
	"log"
	"mailing/utils"
)

func main() {
	// Get the list of valid email addresses
	validEmails, err := utils.GetValidEmails()
	if err != nil {
		log.Fatalf("Error validating emails: %v", err)
	}

	// Send emails only to valid email addresses
	err = utils.SendEmailsToValidRecipients(validEmails)
	if err != nil {
		log.Fatalf("Error sending emails: %v", err)
	}
}
