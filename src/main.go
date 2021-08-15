package main

import (
	"log"
	"net/http"

	"github.com/magiconair/properties"

	"codestep/db"
	"codestep/security"

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

	db.InitConnection(databaseHost, databasePort, databaseDbname, databaseUser, databasePassword)

	mux := http.NewServeMux()

	// Login
	mux.HandleFunc("/login", security.Login)

	// protected multiplexer

	muxProtected := http.NewServeMux()
	muxProtected.HandleFunc("/protected/logout", security.Logout)

	mux.Handle("/protected/", security.ProtectHandler(muxProtected))

	server := &http.Server{
		Addr:    serverAddr + ":" + serverPort,
		Handler: mux,
	}

	log.Fatal(server.ListenAndServeTLS(sslCertificateFilePath, sslPrivateKeyFilePath))

}
