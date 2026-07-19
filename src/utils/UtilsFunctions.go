package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/mail"
	"net/smtp"
	"strconv"

	"github.com/magiconair/properties"
)

var (
	SmtpHost     string
	SmtpPort     string
	SmtpUser     string
	SmtpPassword string
)

func SendEmail(to []string, msg []byte) bool {

	// Set up authentication information.
	auth := smtp.PlainAuth("", SmtpUser, SmtpPassword, SmtpHost)
	log.Print("Smtp auth success")

	// format smtp address
	smtpAddress := fmt.Sprintf("%s:%v", SmtpHost, SmtpPort)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(smtpAddress, auth, SmtpUser, to, msg)
	log.Print("Email sent successfully ")

	if err != nil {
		log.Println(err)
		return false
	}

	// return true on success
	return true
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetClientUrl() (string, string, string) {
	p := properties.MustLoadFile("server.conf", properties.UTF8)
	protocol := p.MustGetString("client_protocol")
	host := p.MustGetString("client_addr")
	port := p.MustGetString("client_port")
	return protocol, host, port
}

func GeneratePassword() string {
	N := rand.Intn(900000) + 100000
	return strconv.Itoa(N)
}
