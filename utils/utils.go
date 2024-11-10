package utils

import (
	"bufio"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"regexp"
)

func SendEmail(to []string, subject, body string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	sender := "vervikar@gmail.com"
	password := "vkeh zuxc ivka kmyp"

	// Формируем сообщение
	msg := "From: " + sender + "\n" +
		"To: " + to[0] + "\n" + // Простой пример, только один получатель
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", sender, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, sender, to, []byte(msg))
	if err != nil {
		return err
	}

	fmt.Println("Email sent successfully")
	return nil
}

func HandleEmails() error {
	file, err := os.Open("db.txt")
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if emailRegex.MatchString(line) {
			fmt.Println("Valid email:", line)

			messageFile, err := os.Open("letter.txt")
			if err != nil {
				return fmt.Errorf("error opening message file: %v", err)
			}
			defer messageFile.Close()

			var body string
			scannerMessage := bufio.NewScanner(messageFile)
			for scannerMessage.Scan() {
				body += scannerMessage.Text() + "\n"
			}

			if err := scannerMessage.Err(); err != nil {
				return fmt.Errorf("error reading message file: %v", err)
			}

			recipients := []string{line}
			subject := "Wisespace"
			err = SendEmail(recipients, subject, body)
			if err != nil {
				log.Printf("Error sending email to %s: %v", line, err)
			}
		} else {
			fmt.Println("Invalid email:", line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	return nil
}
