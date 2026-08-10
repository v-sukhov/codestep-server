package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/magiconair/properties"

	"codestep/db"
	"codestep/liquibase"
	"codestep/security"
	"codestep/services/admin"
	"codestep/services/contest"
	"codestep/services/contest_participation"
	"codestep/services/supertask"
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

	// get dir path
	admin.TmpDirPath = p.MustGetString("tmp_dir")

	// Logging settings
	if p.GetString("cors_log_http", "no") == "yes" {
		CorsLogHttp = true
	} else {
		CorsLogHttp = false
	}

	LogFileName := p.GetString("log_file_name", "")

	/*
		Настройки логирования
	*/
	if LogFileName != "" {
		logFile, err := os.OpenFile(LogFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			log.Println(err)
		} else {
			mw := io.MultiWriter(os.Stderr, logFile)
			log.SetOutput(mw)
		}
	}

	if p.GetString("liquibase_update", "no") == "yes" {
		if err := liquibase.Update(); err != nil {
			log.Fatalf("Liquibase update failed: %v", err)
		}
	}

	db.InitConnection(databaseHost, databasePort, databaseDbname, databaseUser, databasePassword)

	mux := http.NewServeMux()

	// Static files (frontend and static images)

	mux.Handle("/", http.FileServer(http.Dir("./public")))
	mux.Handle("/contests/", &IndexHtmlHandler{})
	mux.Handle("/development/", &IndexHtmlHandler{})
	mux.Handle("/administration/", &IndexHtmlHandler{})

	// Login
	mux.HandleFunc("/api/login", security.Login)

	// Registration
	mux.HandleFunc("/api/register", security.Register)
	mux.HandleFunc("/api/validate-register-code", security.ValidateRegisterCode)
	mux.HandleFunc("/api/resume-register", security.ResumeRegister)

	// Data files (contest logo, supertask logo, supertask images: sprites, backgrounds etc.)
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./files"))))

	// protected multiplexer for api

	muxProtected := http.NewServeMux()

	// Logout
	muxProtected.HandleFunc("/api/protected/logout", security.Logout)
	// Supertask
	muxProtected.HandleFunc("/api/protected/save-supertask", supertask.SaveSupertask)
	muxProtected.HandleFunc("/api/protected/get-supertask", supertask.GetSupertask)
	muxProtected.HandleFunc("/api/protected/get-supertask-versions", supertask.GetSupertaskVersions)
	muxProtected.HandleFunc("/api/protected/get-user-supertask-list", supertask.GetUserSupertaskList)
	// Supertask solution
	muxProtected.HandleFunc("/api/protected/save-supertask-solution", contest_participation.SaveSupertaskSolution)
	muxProtected.HandleFunc("/api/protected/get-supertask-solution", supertask.GetSupertaskSolution)
	// Supertask result
	muxProtected.HandleFunc("/api/protected/save-supertask-result", contest_participation.SaveSupertaskResult)
	muxProtected.HandleFunc("/api/protected/get-supertask-result", supertask.GetSupertaskResult)
	muxProtected.HandleFunc("/api/protected/get-supertask-all-tasks-results", supertask.GetSupertaskAllTasksResults)
	// Contest
	muxProtected.HandleFunc("/api/protected/get-contest", contest.GetContest)
	muxProtected.HandleFunc("/api/protected/get-contest-user-list", contest.GetContestUserList)
	muxProtected.HandleFunc("/api/protected/save-contest", contest.SaveContest)
	muxProtected.HandleFunc("/api/protected/add-supertask-to-contest", contest.AddSupertaskToContest)
	muxProtected.HandleFunc("/api/protected/get-contest-supertask-list", contest_participation.GetContestSupertaskList)
	muxProtected.HandleFunc("/api/protected/get-contest-participation-info", contest_participation.GetContestParticipationInfo)
	muxProtected.HandleFunc("/api/protected/get-user-contest-list", contest.GetUserContestList)
	muxProtected.HandleFunc("/api/protected/remove-supertask-from-contest", contest.RemoveSupertaskFromContest)
	muxProtected.HandleFunc("/api/protected/get-supertask-in-contest-with-results", contest_participation.GetSupertaskInContestWithResults)
	muxProtected.HandleFunc("/api/protected/add-participation-time", contest_participation.AddParticipationTime)
	muxProtected.HandleFunc("/api/protected/edit-user-contest-rights", contest_participation.EditUserContestRights)
	muxProtected.HandleFunc("/api/protected/upload-add-participation-time-file", contest_participation.AddParticipationTimeFile)
	muxProtected.HandleFunc("/api/protected/get-contest-results", contest.GetContestResults)
	muxProtected.HandleFunc("/api/protected/get-contest-results-file", contest.GetContestResultsFile)
	// User managment
	muxProtected.HandleFunc("/api/protected/upload-create-user-list", admin.CreateMultipleUsers)
	muxProtected.HandleFunc("/api/protected/create-contest-users", admin.CreateContestUsers)
	muxProtected.HandleFunc("/api/protected/upload-manage-multiple-users-contest-rights", admin.ManageMultipleUsersContestRights)

	mux.Handle("/api/protected/", security.ProtectHandler(muxProtected))

	server := &http.Server{
		Addr:    serverAddr + ":" + serverPort,
		Handler: CorsHandler(mux),
	}

	log.Fatal(server.ListenAndServeTLS(sslCertificateFilePath, sslPrivateKeyFilePath))

}
