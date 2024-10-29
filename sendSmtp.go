package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/smtp"
	"os"
)

func sendNewListings(msgBody string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	from := "marcus.ikdal2@gmail.com"
	password := os.Getenv("PASSWORD")
	to := []string{"marcus.ikdal@gmail.com"}

	msg := []byte("To: marcus.ikdal@gmail.com\r\n" +
		"Subject: Ny jobbannonse på finn!\r\n" +
		"\r\n" +
		msgBody + "\r\n")

	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")

	return smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
}
