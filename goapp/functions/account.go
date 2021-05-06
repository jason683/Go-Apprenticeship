package functions

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type empty interface{}

//Staff is to be exported
type Staff struct {
	Username string
	Password []byte
	First    string
	Last     string
	Email    string
	Rights   string
}

//UserResult is to be exported
type UserResult struct {
	Username string
}

var (
	//Db is to be exported
	Db *sql.DB
	//Err is to be exported
	Err         error
	mapUsers    = map[string]Staff{}
	mapSessions = map[string]string{}
	mapBusiness = map[string]string{}
	relationMap = map[string]string{}
	//Tpl is to be exported
	Tpl *template.Template
)

//Start is to be exported
func Start(res http.ResponseWriter, req *http.Request) {
	//this will retrieve the Staff information and display it on the webpage
	myUser := GetUser(res, req)
	Tpl.ExecuteTemplate(res, "index.html", myUser)
}

//GetUser is to be exported
func GetUser(res http.ResponseWriter, req *http.Request) Staff {
	//this stores the cookie data in myCookie
	//this subsequent part is necessary so that when the user accesses the index page, the page can load
	//without this section, the program will panic
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		id := uuid.NewV4()
		myCookie = &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
			MaxAge:   2600000,
		}
		http.SetCookie(res, myCookie)
	}
	log.Println("Cookie session has been created/validated")
	//this will check if there is a username already tagged to mapSessions, and
	//if yes, the information details of the user will be stored in a new variable called myUser
	//this is solely for the start function to know the user details so that the correct information can be displayed
	var myUser Staff
	if username, ok := mapSessions[myCookie.Value]; ok {
		myUser = mapUsers[username]
		log.Println("user details are now stored in myUser")
	}
	return myUser
}

//Signup is to be exported
func Signup(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	var myUser Staff
	//request for user input
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		email := req.FormValue("email")

		//input validation
		errorMessage := map[string]string{}
		if username == "" || password == "" || firstname == "" || lastname == "" || email == "" {
			log.Println("There are missing particulars")
			errorMessage["input1"] = "Missing particulars"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input1")
			return
		}
		//this Dbquery command will return all data from the entire Users table
		results, err := Db.Query("SELECT Username, Password, FirstName, LastName, Email from Users")
		if err != nil {
			panic(err.Error())
		}
		//check if there is already an existing user
		checkUsers := map[string]Staff{}
		var existingUser Staff
		//analysing row by row
		for results.Next() {
			err = results.Scan(&existingUser.Username, &existingUser.Password, &existingUser.First, &existingUser.Last, &existingUser.Email)
			if err != nil {
				panic(err.Error())
			}
			//store all the Db information in checkUsers
			checkUsers[existingUser.Username] = existingUser
		}
		for k := range checkUsers {
			if username == k {
				log.Println("Username has been taken")
				errorMessage["input2"] = "Username has been taken"
				Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
				delete(errorMessage, "input2")
			}
		}
		//empty out checkUsers since not needed anymore
		delete(checkUsers, existingUser.Username)
		//input validation for username and password
		reg, err := regexp.Compile("^[a-zA-Z0-9]*$")
		if err != nil {
			log.Fatal(err)
		}
		if reg.MatchString(username) == false {
			log.Println("username contains non-alphanumeric characters")
			errorMessage["input3"] = "username can only have alphanumeric characters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input3")
			return
		}
		if reg.MatchString(password) == false {
			log.Println("password contains non-alphanumeric characters")
			errorMessage["input4"] = "password can only have alphanumeric characters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input4")
			return
		}
		if len(username) < 5 || len(username) > 12 {
			log.Println("username contains too few or too many characters")
			errorMessage["input5"] = "username must have at least 5 characters and at most 12 characters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input5")
			return
		}
		//password must have at least 8 characters
		if len(password) < 8 || len(password) > 15 {
			log.Println("password contains too few or too many characters")
			errorMessage["input6"] = "password must have at least 8 characters and at most 15 characters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input6")
			return
		}
		//this regex expression ensures first and last name are both entered correctly
		rege, err := regexp.Compile("^[a-zA-Z]*$")
		if err != nil {
			log.Fatal(err)
		}
		if rege.MatchString(firstname) == false {
			log.Println("First name contains non-alphabetical letters")
			errorMessage["input7"] = "first name can only have alphabetical letters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input7")
			return
		}
		if rege.MatchString(lastname) == false {
			log.Println("Last name contains non-alphabetical letters")
			errorMessage["input8"] = "last name can only have alphabetical letters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input8")
			return
		}
		//first name can have at most 20 characters
		if len(firstname) > 20 {
			log.Println("First name has too many letters")
			errorMessage["input9"] = "first name can have at most 20 characters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input9")
			return
		}
		//last name can have at most 20 characters
		if len(lastname) > 20 {
			log.Println("Last name has too many letters")
			errorMessage["input10"] = "last name can have at most 20 characters"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input10")
			return
		}
		//this expression ensures that the email is entered correctly
		regee, err := regexp.Compile("^[a-zA-Z0-9_.+]+@[a-zA-Z0-9]+.[a-zA-Z0-9]+$")
		if err != nil {
			log.Fatal(err)
		}
		if regee.MatchString(email) == false {
			log.Println("email format is wrong")
			errorMessage["input11"] = "email format is wrong"
			Tpl.ExecuteTemplate(res, "signup.html", errorMessage)
			delete(errorMessage, "input11")
			return
		}
		//defer function in response to potential panic
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered")
			}
		}()
		//using bcrypt package to encrypt the password
		bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(res, "Internal server error", http.StatusInternalServerError)
		}
		//if inputs have all been validated, store the data in myUser variable
		myUser = Staff{username, bPassword, firstname, lastname, email, ""}
		//then create a new key value pair in mapUsers (key is username and values are the rest of the information)
		mapUsers[username] = myUser
		//sprintf function prints to the string
		query := fmt.Sprintf("INSERT INTO Users VALUES ('%s', '%s', '%s', '%s', '%s', NULL)", myUser.Username, myUser.Password, myUser.First, myUser.Last, myUser.Email)
		//query once again to check if there are any errors in this process
		_, err = Db.Query(query)
		if err != nil {
			panic(err.Error())
		}
		//Set cookie process
		id := uuid.NewV4()
		expireTime := time.Now().Add(30 * time.Minute)
		myCookie := &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  expireTime,
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
		}
		http.SetCookie(res, myCookie)
		log.Println("New cookie session created")
		//create a new key value pair in mapSessions to track session usage
		mapSessions[myCookie.Value] = username
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
		return
	}
	//if the user clicks on this page, he or she will first see this template
	Tpl.ExecuteTemplate(res, "signup.html", nil)
}

//AlreadyLoggedIn is to be exported
func AlreadyLoggedIn(req *http.Request) bool {
	//this will request for the http cookie and if there is one already (i.e. the user has an existing session),
	//the function will return true
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	username := mapSessions[myCookie.Value]
	_, ok := mapUsers[username]
	return ok
}

//Login is to be exported
func Login(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		results, err := Db.Query("SELECT Username, Password, FirstName, LastName, Email, Rights FROM contracts_db.Users WHERE Username = ?", username)
		if err != nil {
			panic(err.Error())
		}
		errorMessage := map[string]string{}
		var myUser Staff
		for results.Next() {
			err := results.Scan(&myUser.Username, &myUser.Password, &myUser.First, &myUser.Last, &myUser.Email, &myUser.Rights)
			if err != nil {
				errorMessage["errorUserAndPassword"] = "Username and/or password do not match or are invalid"
				Tpl.ExecuteTemplate(res, "login.html", errorMessage)
				delete(errorMessage, "errorUserAndPassword")
				return
			}
		}
		//using bcrypt library to check if the encrypted passwords match
		err = bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			errorMessage["errorUserAndPassword"] = "Username and/or password do not match or are invalid"
			Tpl.ExecuteTemplate(res, "login.html", errorMessage)
			delete(errorMessage, "errorUserAndPassword")
			return
		}
		//create a new cookie session once the user manages to login
		id := uuid.NewV4()
		expireTime := time.Now().Add(30 * time.Minute)
		myCookie := &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  expireTime,
			HttpOnly: true,
			Path:     "/",
			Domain:   "127.0.0.1",
			Secure:   true,
		}
		http.SetCookie(res, myCookie)
		//create new key value pairs for mapSessions and mapUsers
		mapSessions[myCookie.Value] = username
		mapUsers[username] = myUser
		log.Println("Cookie session created successfully")
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
		return
	}
	Tpl.ExecuteTemplate(res, "login.html", nil)
}

//Logout is to be exported
func Logout(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	//deleting the key value pairs in mapUsers and mapSessions. Order is important
	myCookie, _ := req.Cookie("myCookie")
	username := mapSessions[myCookie.Value]
	delete(mapSessions, myCookie.Value)
	//this will delete the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	delete(mapUsers, username)
	log.Println("Cookie successfully deleted")
	Tpl.ExecuteTemplate(res, "logout.html", nil)
}
