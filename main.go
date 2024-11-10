package main

import (
	"log"
	"mailing/utils"
)

func main() {
	err := utils.HandleEmails()
	if err != nil {
		log.Fatalf("Error handling emails: %v", err)
	}
}
