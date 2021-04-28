package functions

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//CreateRequest to be exported
func CreateRequest(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser := GetUser(res, req)
	if req.Method == http.MethodPost {
		errorMessage := map[string]string{}

		signingEntity := req.FormValue("signingentity")
		counterpartyName := req.FormValue("counterpartyname")
		business := req.FormValue("business")
		businessOwner := req.FormValue("businessowner")
		approveStatus := "Pending"
		financeTax := "Pending"
		signed := "Pending"
		if signingEntity == "" || counterpartyName == "" || business == "" || businessOwner == "" {
			errorMessage["input1"] = "Did you miss out entering any fields?"
			Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
			delete(errorMessage, "input1")
			return
		}
		contractValue := req.FormValue("contractvalue")
		_, err := strconv.Atoi(contractValue)
		if err != nil {
			errorMessage["input0"] = "Contract value has to be an integer"
			Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
			delete(errorMessage, "input0")
			return
		}

		timeAction := time.Now().Format(time.RFC3339)
		//NULL can be used to circumvent the int auto increment in sql
		query := fmt.Sprintf("INSERT INTO Contracts VALUES (NULL, '%s', '%s', '%s', '%s', '%s', '%s', NULL, '%s', '%s', NULL, '%s', NULL, '%s', NULL)", signingEntity, counterpartyName, business, myUser.Username, businessOwner, approveStatus, financeTax, contractValue, signed, timeAction)
		_, err = Db.Query(query)
		if err != nil {
			fmt.Println(err)
		}

		SendEmail("testtechnology.93@gmail.com")
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
		return
	}
	Tpl.ExecuteTemplate(res, "requestform.html", nil)
}
