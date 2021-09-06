package utils

import (
	"fmt"
	"log"
    "net/mail"
	"net/smtp"
	"github.com/magiconair/properties"
)


func SendEmail(to []string, msg []byte) bool {
	p := properties.MustLoadFile("server.conf", properties.UTF8)

    // email properties
	smtpAddr := p.MustGetString("smtp_host")
	smtpPort := p.MustGetString("smtp_port")
	smptUser := p.MustGetString("smtp_user")
	smtpPass := p.MustGetString("smtp_password")

    // Set up authentication information.
    auth := smtp.PlainAuth("", smptUser, smtpPass, smtpAddr)

    // format smtp address
    smtpAddress := fmt.Sprintf("%s:%v", smtpAddr, smtpPort)

    // Connect to the server, authenticate, set the sender and recipient,
    // and send the email all in one step.
    err := smtp.SendMail(smtpAddress, auth, smptUser, to, msg)

    if err != nil {
        log.Fatal(err)
        return false
    }

    // return true on success
    return true
}

func GetRegisterSecret() string {
	p := properties.MustLoadFile("server.conf", properties.UTF8)
	secretJwt := p.MustGetString("jwt_secret")
	return secretJwt
}

func ValidEmail(email string) bool {
    _, err := mail.ParseAddress(email)
    return err == nil
}