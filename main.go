package main

import (
	"log"
	"mailing/utils"
)

func main() {
	// Запуск обработки email-адресов
	err := utils.HandleEmails()
	if err != nil {
		log.Fatalf("Error handling emails: %v", err)
	}
}
