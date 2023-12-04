package main

import (
	"log"
	"net/http"

	"github.com/magiconair/properties"

	"codestep/db"
	"codestep/security"
	"codestep/services"
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
	utils.SmtpPassword = p.MustGetString("smtp_password")

	// register settings
	security.JwtSecret = p.MustGetString("jwt_secret")
	security.JwtRegisterSecret = p.MustGetString("jwt_register_secret")
	security.JwtTokenLifetimeMinute = p.MustGetInt("jwt_token_lifetime_minute")
	security.JwtRegisterTokenLifetimeMinute = p.MustGetInt("jwt_register_token_lifetime_minute")

	// Logging settings
	if p.GetString("cors_log_http", "no") == "yes" {
		CorsLogHttp = true
	} else {
		CorsLogHttp = false
	}

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
	// Logout
	muxProtected.HandleFunc("/api/protected/logout", security.Logout)
	// Supertask
	muxProtected.HandleFunc("/api/protected/save-supertask", services.SaveSupertask)
	muxProtected.HandleFunc("/api/protected/get-supertask", services.GetSupertask)
	muxProtected.HandleFunc("/api/protected/get-user-supertask-list", services.GetUserSupertaskList)
	muxProtected.HandleFunc("/api/protected/save-supertask-solution", services.SaveSupertaskSolution)
	muxProtected.HandleFunc("/api/protected/get-supertask-solution", services.GetSupertaskSolution)
	muxProtected.HandleFunc("/api/protected/save-supertask-result", services.SaveSupertaskResult)

	mux.Handle("/api/protected/", security.ProtectHandler(muxProtected))

	server := &http.Server{
		Addr:    serverAddr + ":" + serverPort,
		Handler: CorsHandler(mux),
	}

	log.Fatal(server.ListenAndServeTLS(sslCertificateFilePath, sslPrivateKeyFilePath))

}
