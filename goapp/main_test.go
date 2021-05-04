package main

import (
	"Apprentice/Go-Apprenticeship/goapp/functions"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type account struct {
	username  string
	bPassword []byte
}

func Router() *mux.Router {
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
	router.HandleFunc("/result", functions.Result)
	return router
}

func TestStart(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)

	//assert.Equal(t, 200, response.Code, "OK response is expected")
	res := response.Result()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}
	defer res.Body.Close()
	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("could not read response: %v", err)
	}
}

func TestLogin(t *testing.T) {
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
	testUser := "jasonfoo1"
	testPassword := "hello123"
	results, err := functions.Db.Query("SELECT Username, Password FROM contracts_db.Users WHERE Username = ?", testUser)
	if err != nil {
		fmt.Println(err)
	}
	var myUser account
	for results.Next() {
		err := results.Scan(&myUser.username, &myUser.bPassword)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = bcrypt.CompareHashAndPassword(myUser.bPassword, []byte(testPassword))
	if err != nil {
		t.Errorf("Expect nil error; got %v", err)
	}
	requestBody, err := json.Marshal(map[string]string{
		"username": testUser,
		"password": testPassword,
	})
	if err != nil {
		log.Fatalln(err)
	}

	request, err := http.NewRequest("POST", "/login", bytes.NewReader(requestBody))

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)
	//functions.Login(response, request)

	res := response.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	fmt.Println(string(b))
	templateString := string(b)
	boolean0 := strings.Contains(templateString, "Show/Upload contracts")
	if boolean0 == true {
		t.Errorf("Expected false; got %v", boolean0)
	}
	boolean1 := strings.Contains(templateString, "Create a request")
	if boolean1 == false {
		t.Errorf("Expected true; got %v", boolean1)
	}
}

// func TestSignup(t *testing.T) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	databaseUser := os.Getenv("DATABASEUSER")
// 	databasePW := os.Getenv("DATABASEPW")
// 	databaseString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/contracts_db?parseTime=true", databaseUser, databasePW)
// 	functions.Db, functions.Err = sql.Open("mysql", databaseString)
// 	if functions.Err != nil {
// 		panic(functions.Err.Error())
// 	} else {
// 		fmt.Println("Database opened!")
// 	}
// 	defer functions.Db.Close()
// 	testUser := struct {
// 		username  string
// 		password  string
// 		firstName string
// 		lastName  string
// 		email     string
// 	}{
// 		username: "test", password: "testing", firstName: "first", lastName: "last", email: "abc@gmail.com",
// 	}
// 	request, _ := http.NewRequest("POST", "https://127.0.0.1:5000/signup?username="+testUser.username+"&password="+testUser.password+"&firstname="+testUser.firstName+"&lastname="+testUser.lastName+"&email="+testUser.email, nil)
// 	response := httptest.NewRecorder()
// 	Router().ServeHTTP(response, request)

// 	res := response.Result()
// 	defer res.Body.Close()
// 	b, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatalf("Error: %v", err)
// 	}
// 	templateString := string(b)
// 	boolean := strings.Contains(templateString, "Missing particulars")
// 	if boolean == false {
// 		t.Errorf("Expected true; got %v", boolean)
// 	}
// }

// func TestDirectory(t *testing.T) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	databaseUser := os.Getenv("DATABASEUSER")
// 	databasePW := os.Getenv("DATABASEPW")
// 	databaseString := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/contracts_db?parseTime=true", databaseUser, databasePW)
// 	functions.Db, functions.Err = sql.Open("mysql", databaseString)
// 	if functions.Err != nil {
// 		panic(functions.Err.Error())
// 	} else {
// 		fmt.Println("Database opened!")
// 	}
// 	defer functions.Db.Close()
// 	testUser := "jasonfoo1"
// 	testPassword := "hello123"
// 	request, _ := http.NewRequest("POST", "https://127.0.0.1:5000/login?username="+testUser+"&password="+testPassword, nil)
// 	fmt.Println(request)
// 	fmt.Println()

// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

// 	fmt.Println(request.Body)
// 	fmt.Println()
// 	response := httptest.NewRecorder()

// 	Router().ServeHTTP(response, request)
// 	res := response.Result()
// 	b, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		t.Fatalf("Error: %v", err)
// 	}
// 	templateString := string(b)
// 	boolean0 := strings.Contains(templateString, "Show/Upload contracts")
// 	if boolean0 == true {
// 		t.Errorf("Expected false; got %v", boolean0)
// 	}
// 	boolean1 := strings.Contains(templateString, "Create a request")
// 	if boolean1 == false {
// 		t.Errorf("Expected true; got %v", boolean1)
// 	}
// }
