package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func financeTax(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := getUser(res, req)
	if myUser.Rights == "financetax" {
		results, err := db.Query("SELECT Id, SigningEntity, CounterpartyName, Business, Requester FROM contracts_db.Contracts WHERE FinanceTax = 'Pending'")
		if err != nil {
			panic(err.Error())
		}
		//display variable will contain a list of all the pending contract requests
		display := []contractRequest{}
		var reviewRequest contractRequest
		for results.Next() {
			err := results.Scan(&reviewRequest.ID, &reviewRequest.SigningEntity, &reviewRequest.CounterpartyName, &reviewRequest.Business, &reviewRequest.Requester)
			if err != nil {
				panic(err.Error())
			}
			display = append(display, reviewRequest)
		}
		if req.Method == http.MethodPost {
			if !alreadyLoggedIn(req) {
				http.Redirect(res, req, "/", http.StatusSeeOther)
			}
			contractRequestIDstring := req.FormValue("contractrequestid")
			contractRequestIDint, err := strconv.Atoi(contractRequestIDstring)
			if err != nil {
				panic(err.Error())
			}
			//section below will change pending status to either approve or reject status
			contractRequestStatus := req.FormValue("approvereject")
			for _, v := range display {
				if v.ID == contractRequestIDint {
					lowercaseContractRequestStatus := strings.ToLower(contractRequestStatus)
					if lowercaseContractRequestStatus == "approve" || lowercaseContractRequestStatus == "reject" {
						query := fmt.Sprintf("UPDATE Contracts SET FinanceTax='%s' WHERE Id='%s'", contractRequestStatus, contractRequestIDstring)
						_, err := db.Query(query)
						if err != nil {
							panic(err.Error())
						}
					}
				}
			}
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
		}
		tpl.ExecuteTemplate(res, "financetax.html", display)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
