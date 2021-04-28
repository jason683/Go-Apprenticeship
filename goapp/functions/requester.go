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
	if myUser.Rights == "bizrequester" {
		mapBusiness["first"] = "Games & E-Commerce"
		mapBusiness["second"] = "Digital Finance"

		if req.Method == http.MethodPost {
			errorMessage := map[string]string{}

			signingEntity := req.FormValue("signingentity")
			counterpartyName := req.FormValue("counterpartyname")
			business := req.FormValue("business")
			contractType := req.FormValue("contracttype")
			businessOwner := req.FormValue("businessowner")
			newBusiness := req.FormValue("newbusiness")
			approveStatus := "Pending"
			financeTax := "Pending"
			signed := "Pending"
			if signingEntity == "" || counterpartyName == "" || contractType == "" || businessOwner == "" {
				errorMessage["input1"] = "Did you miss out entering any fields?"
				Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
				delete(errorMessage, "input1")
				return
			}
			if business == "" && newBusiness == "" {
				errorMessage["nobusiness"] = "You have not entered any value for the business field"
				Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
				delete(errorMessage, "nobusiness")
				return
			}
			if newBusiness != "" && business != "" {
				errorMessage["duplicatebusiness"] = "You have entered values for both business fields"
				Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
				delete(errorMessage, "duplicatebusiness")
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
			query := fmt.Sprintf("INSERT INTO Contracts VALUES (NULL, '%s', '%s', '%s', '%s', '%s', '%s', '%s', NULL, '%s', '%s', NULL, '%s', NULL, '%s', NULL)", signingEntity, counterpartyName, business, contractType, myUser.Username, businessOwner, approveStatus, financeTax, contractValue, signed, timeAction)
			_, err = Db.Query(query)
			if err != nil {
				fmt.Println(err)
			}

			SendEmail("testtechnology.93@gmail.com")
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
			return
		}
		Tpl.ExecuteTemplate(res, "requestform.html", mapBusiness)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
