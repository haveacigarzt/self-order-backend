package main

import (
	"fmt"
	"net/smtp"
)

func main() {
	auth := smtp.PlainAuth(
		"",
		"pramanaafriandy@gmail.com",
		"ypcujvhuvnebgmuu",
		"smtp.gmail.com",
	)

	msg := "Subject: My subject\nThis is the body of my email"

	err := smtp.SendMail("smtp.gmail.com:587", auth, "pramanaafriandy@gmail.com", []string{"afriandy193@gmail.com"}, []byte(msg))

	if err != nil {
		fmt.Println(err)
	}
}
