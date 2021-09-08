package main

import (
	"log"
	"net/http"

	"github.com/magiconair/properties"

	"codestep/db"
	"codestep/security"
	"codestep/utils"

	_ "codestep/docs"
)

func main() {
	p := properties.MustLoadFile("server.conf", properties.UTF8)

	// server properties
	serverAddr := p.MustGetString("server_addr")
	serverPort := p.MustGetString("server_port")

	// database connection properties
	databaseHost := p.MustGetString("database_host")
	databasePort := p.MustGetString("database_port")
	databaseDbname := p.MustGetString("database_dbname")
	databaseUser := p.MustGetString("database_user")
	databasePassword := p.MustGetString("database_password")

	// password encryption
	db.EncryptionSaltWord = p.MustGetString("encryption_salt_word")

	// certificate files
	sslCertificateFilePath := p.MustGetString("ssl_certificate")
	sslPrivateKeyFilePath := p.MustGetString("ssl_private_key")

	// SMTP settings
	utils.SmtpHost = p.MustGetString("smtp_host")
	utils.SmtpPort = p.MustGetString("smtp_port")
	utils.SmtpUser = p.MustGetString("smtp_user")
	utils.SmtpPass = p.MustGetString("smtp_password")

	// register settings
	security.JwtSecrete = p.MustGetString("jwt_secret")
	security.JwtTokenLifetimeMinute = p.MustGetInt("jwt_token_lifetime_minute")

	db.InitConnection(databaseHost, databasePort, databaseDbname, databaseUser, databasePassword)

	mux := http.NewServeMux()

	// Login
	mux.HandleFunc("/api/login", security.Login)

	// Registration
	mux.HandleFunc("/api/register", security.Register)
	mux.HandleFunc("/api/validate-register-code", security.ValidateRegisterCode)
	mux.HandleFunc("/api/resume-register", security.ResumeRegister)

	// protected multiplexer

	muxProtected := http.NewServeMux()
	muxProtected.HandleFunc("/api/protected/logout", security.Logout)

	mux.Handle("/api/protected/", security.ProtectHandler(muxProtected))

	server := &http.Server{
		Addr:    serverAddr + ":" + serverPort,
		Handler: CorsHandler(mux),
	}

	log.Fatal(server.ListenAndServeTLS(sslCertificateFilePath, sslPrivateKeyFilePath))

}
