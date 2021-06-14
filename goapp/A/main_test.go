package main

import (
	"Apprentice/Go-Apprenticeship/goapp/functions"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bou.ke/monkey"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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

	request, err := http.NewRequest("POST", fmt.Sprintf("https://127.0.0.1:5000/login?username=%s&password=%s", testUser, testPassword), nil)
	if err != nil {
		t.Error(err)
	}

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)

	res := response.Result()
	defer res.Body.Close()
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	boolean0 := assert.Equal(t, "/directory", res.Header.Get("Location"))
	if boolean0 != true {
		t.Error(err)
	}
}

func TestSignup(t *testing.T) {
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
	testUser := []struct {
		username  string
		password  string
		firstName string
		lastName  string
		email     string
	}{
		{username: "test", password: "testing", firstName: "first", lastName: "last", email: "abc@gmail.com"},
		{username: "test1", password: "testing", email: "abc@gmail.com"},
		{username: "test2", password: "testing1", firstName: "first", lastName: "last", email: "abc@gmail.com"},
	}

	errorCount := 0
	redirectionCount := 0
	for _, tc := range testUser {
		t.Run(tc.username, func(t *testing.T) {
			request, _ := http.NewRequest("POST", "https://127.0.0.1:5000/signup?username="+tc.username+"&password="+tc.password+"&firstname="+tc.firstName+"&lastname="+tc.lastName+"&email="+tc.email, nil)
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)

			res := response.Result()
			defer res.Body.Close()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			templateString := string(b)

			boolean0 := checkSubstrings(templateString, "Missing particulars", "username must have at least 5 characters and at most 12 characters")
			if boolean0 == true {
				errorCount++
			}

			boolean1 := assert.NotEqual(t, "/directory", res.Header.Get("Location"))
			if boolean1 == false {
				redirectionCount++
			}
		})
	}
	if errorCount != 2 {
		t.Errorf("Expected exactly 2 cases to generate an error each; got %v cases", errorCount)
	}
	if redirectionCount != 0 {
		t.Errorf("Expected exactly 0 redirection case; got %v", redirectionCount)
	}
}

func checkSubstrings(template string, substrings ...string) bool {
	for _, substring := range substrings {
		if strings.Contains(template, substring) {
			return true
		}
	}
	return false
}

func TestDirectory(t *testing.T) {
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
	request, err := http.NewRequest("GET", fmt.Sprintf("https://127.0.0.1:5000/directory"), nil)
	if err != nil {
		t.Error(err)
	}
	response := httptest.NewRecorder()

	g := monkey.Patch(functions.AlreadyLoggedIn, func(req *http.Request) bool {
		return true
	})
	defer g.Unpatch()

	h := monkey.Patch(functions.GetUser, func(res http.ResponseWriter, req *http.Request) functions.Staff {
		return functions.Staff{Rights: "bizrequester"}
	})
	defer h.Unpatch()

	Router().ServeHTTP(response, request)
	res := response.Result()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	templateString := string(b)
	fmt.Println(templateString)
	boolean0 := strings.Contains(templateString, "Show/Upload contracts")
	if boolean0 == true {
		t.Errorf("Expected false; got %v", boolean0)
	}
	fmt.Println(templateString)

	boolean1 := strings.Contains(templateString, "Create a request")
	if boolean1 == false {
		t.Errorf("Expected true; got %v", boolean1)
	}
}

func TestBizRequest(t *testing.T) {
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
	testUser := []struct {
		SigningEntity    string
		CounterpartyName string
		BusinessType     string
		ContractType     string
		ContractValue    string
		BusinessOwner    string
		EffectiveDate    string
		TerminationDate  string
	}{
		{SigningEntity: "Shopee", CounterpartyName: "Roblox", BusinessType: "Hardware"},
		{SigningEntity: "Shopee", CounterpartyName: "Facebook", BusinessType: "Digital Finance", ContractType: "Tech", ContractValue: "Hey", BusinessOwner: "testresults", EffectiveDate: "2019-01-01", TerminationDate: "2020-01-01"},
		{SigningEntity: "Shopee", CounterpartyName: "Facebook", BusinessType: "testresults", ContractType: "Tech", ContractValue: "Hey", BusinessOwner: "testresults", EffectiveDate: "2019-01-01", TerminationDate: "2020-01-01"},
	}

	errorCount := 0
	redirectionCount := 0
	for _, tc := range testUser {
		t.Run(tc.SigningEntity, func(t *testing.T) {

			g := monkey.Patch(functions.AlreadyLoggedIn, func(req *http.Request) bool {
				return true
			})
			defer g.Unpatch()

			h := monkey.Patch(functions.GetUser, func(res http.ResponseWriter, req *http.Request) functions.Staff {
				return functions.Staff{Rights: "bizrequester"}
			})
			defer h.Unpatch()
			request, _ := http.NewRequest("POST", "https://127.0.0.1:5000/createrequest?signingentity="+tc.SigningEntity+"&counterpartyname="+tc.CounterpartyName+"&business="+tc.BusinessType+"&contracttype="+tc.ContractType+"&contractvalue="+tc.ContractValue+"&businessowner="+tc.BusinessOwner+"&effectivedate="+tc.EffectiveDate+"&terminationdate="+tc.TerminationDate, nil)
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)

			res := response.Result()
			defer res.Body.Close()
			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}
			templateString := string(b)
			boolean0 := checkSubstrings(templateString, "Contract value has to be an integer", "Did you miss out entering any of the compulsory fields?", "You need to key in a valid business type")
			if boolean0 == true {
				errorCount++
			}
			boolean1 := assert.NotEqual(t, "/directory", res.Header.Get("Location"))
			if boolean1 == false {
				redirectionCount++
			}
		})
	}
	if errorCount != 3 {
		t.Errorf("Expected all 3 cases to generate an error each; got %v cases", errorCount)
	}
	if redirectionCount != 0 {
		t.Errorf("Expected exactly 0 redirection case; got %v", redirectionCount)
	}
}

func TestUpload(t *testing.T) {
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

	g := monkey.Patch(functions.AlreadyLoggedIn, func(req *http.Request) bool {
		return true
	})
	defer g.Unpatch()

	h := monkey.Patch(functions.GetUser, func(res http.ResponseWriter, req *http.Request) functions.Staff {
		return functions.Staff{Rights: "legal"}
	})
	defer h.Unpatch()

	file, err := os.Open("test.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("myFile", filepath.Base(file.Name()))

	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", "https://127.0.0.1:5000/upload", body)

	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()

	Router().ServeHTTP(response, request)

	res := response.Result()
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	templateString := string(b)
	fmt.Println(templateString)

}
