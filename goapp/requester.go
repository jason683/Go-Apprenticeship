package main

import (
	"fmt"
	"net/http"
)

func createRequest(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := getUser(res, req)
	if req.Method == http.MethodPost {
		signingEntity := req.FormValue("signingentity")
		counterpartyName := req.FormValue("counterpartyname")
		business := req.FormValue("business")
		businessOwner := req.FormValue("businessowner")
		approveStatus := "Pending"
		financeTax := "Pending"
		errorMessage := map[string]string{}
		if signingEntity == "" || counterpartyName == "" || business == "" || businessOwner == "" {
			errorMessage["input1"] = "Did you miss out entering any fields?"
			tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
			delete(errorMessage, "input1")
			return
		}
		//NULL can be used to circumvent the int auto increment in sql
		query := fmt.Sprintf("INSERT INTO Contracts VALUES (NULL, '%s', '%s', '%s', '%s', '%s', '%s', NULL, '%s')", signingEntity, counterpartyName, business, myUser.Username, businessOwner, approveStatus, financeTax)
		fmt.Println("test")
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Hello world")
		}
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(res, "requestform.html", nil)
}
