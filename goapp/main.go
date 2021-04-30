package main

import (
	"Apprentice/Go-Apprenticeship/goapp/functions"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	functions.Tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	databaseUser := os.Getenv("DATABASEUSER")
	databasePW := os.Getenv("DATABASEPW")
	databaseString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/contracts_db?parseTime=true", databaseUser, databasePW)
	functions.Db, functions.Err = sql.Open("mysql", databaseString)
	if functions.Err != nil {
		panic(functions.Err.Error())
	} else {
		fmt.Println("Database opened!")
	}
	defer functions.Db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/", functions.Start)
	router.HandleFunc("/signup", functions.Signup)
	router.HandleFunc("/login", functions.Login)
	router.HandleFunc("/logout", functions.Logout)
	router.HandleFunc("/directory", functions.Directory)
	router.HandleFunc("/createrequest", functions.CreateRequest)
	router.HandleFunc("/reviewrequest", functions.ReviewRequest)
	router.HandleFunc("/fintax", functions.FinanceTax)
	router.HandleFunc("/upload", functions.UploadFile)
	router.HandleFunc("/reviewcontractvalue", functions.ValueApproval)
	router.HandleFunc("/showcontracts", functions.ShowContracts)
	router.HandleFunc("/archive", functions.ArchiveContract)
	router.HandleFunc("/outdated", functions.IdentifyOutdatedRequest)
	router.HandleFunc("/emaillist", functions.EmailList)
	router.HandleFunc("/readfile", functions.Test)
	fmt.Println("Listening at port 5000")

	http.ListenAndServeTLS(":5000", "cert.pem", "key.pem", router)

}
