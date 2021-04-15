package main

import (
	"fmt"
	"net/http"
)

func showContracts(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := getUser(res, req)
	if myUser.Rights == "signatory" {
		//to add on new column (signed?)
		results, err := db.Query("SELECT Contract FROM contracts_db.Contracts WHERE SeniorFinance ='yes' AND Contract IS NOT NULL")
		if err != nil {
			fmt.Println("Cannot extract contract file")
		}
		//display is a list of strings
		display := []string{}
		var row string
		for results.Next() {
			err := results.Scan(&row)
			if err != nil {
				fmt.Println("Cannot scan into row")
			}
			display = append(display, row)
		}

		tpl.ExecuteTemplate(res, "signatory.html", display)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
