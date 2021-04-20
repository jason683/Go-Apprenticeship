package main

import (
	"Apprentice/Go-Apprenticeship/goapp/functions"
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func init() {
	functions.Tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	functions.Db, functions.Err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/contracts_db")
	// var databaseUser string
	// databaseUser = os.Getenv("DATABASEUSER")
	// var databasePW string
	// databasePW = os.Getenv("DATABASEPW")
	// databaseString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/items_db", databaseUser, databasePW)
	// db, err = sql.Open("mysql", databaseString)
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
	fmt.Println("Listening at port 5000")
	http.ListenAndServeTLS(":5000", "cert.pem", "key.pem", router)
}
