package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type empty interface{}

type staff struct {
	Username string
	Password []byte
	First    string
	Last     string
	Email    string
	Rights   string
}

var (
	db            *sql.DB
	err           error
	mapUsers      = map[string]staff{}
	mapSessions   = map[string]string{}
	mapTypeRights = map[string]string{}
	tpl           *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	var databaseUser string
	databaseUser = os.Getenv("DATABASEUSER")
	var databasePW string
	databasePW = os.Getenv("DATABASEPW")
	databaseString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/items_db", databaseUser, databasePW)
	db, err = sql.Open("mysql", databaseString)
	if err != nil {
		panic(err.Error())
	} else {
		fmt.Println("Database opened!")
	}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/", start)
	router.HandleFunc("/signup", signup)
	router.HandleFunc("/login", login)
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/directory", directory)
	router.HandleFunc("/createrequest", createRequest)
	router.HandleFunc("/reviewrequest", reviewRequest)
	router.HandleFunc("/fintax", financeTax)
	router.HandleFunc("/upload", uploadFile)
	router.HandleFunc("/reviewcontractvalue", valueApproval)
	router.HandleFunc("/showcontracts", showContracts)
	fmt.Println("Listening at port 5000")
	http.ListenAndServeTLS(":5000", "cert.pem", "key.pem", router)
}
