package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type valueRequest struct {
	ID               int
	SigningEntity    string
	CounterpartyName string
	Business         string
	Value            string
}

func valueApproval(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := getUser(res, req)
	if myUser.Rights == "contractadmin" {
		results, err := db.Query("SELECT Id, SigningEntity, CounterpartyName, Business, ContractValue FROM contracts_db.Contracts WHERE ContractValue IS NOT NULL AND FinanceTax = 'Pending' AND SeniorFinance IS NULL")
		if err != nil {
			fmt.Println("Something has happened")
		}
		//display variable will contain a list of all unapproved contract requests in relation to their contract value
		display := []valueRequest{}
		var reviewRequest valueRequest
		for results.Next() {
			err := results.Scan(&reviewRequest.ID, &reviewRequest.SigningEntity, &reviewRequest.CounterpartyName, &reviewRequest.Business, &reviewRequest.Value)
			if err != nil {
				fmt.Println("Unable to scan into reviewRequest variable")
			}
			display = append(display, reviewRequest)
		}

		if req.Method == http.MethodPost {
			contractRequestIDstring := req.FormValue("contractrequestid")
			contractRequestIDint, err := strconv.Atoi(contractRequestIDstring)
			if err != nil {
				fmt.Println("invalid, unapproved or no contract ID")
			}

			//change NULL to potentially yes. Yes if contract value is beyond a certain amount for the given business
			signatory := req.FormValue("signatory")
			for _, v := range display {
				if v.ID == contractRequestIDint {
					lowercaseSignatory := strings.ToLower(signatory)
					if lowercaseSignatory == "yes" {
						query := fmt.Sprintf("UPDATE Contracts SET SeniorFinance='%s' WHERE Id='%s'", signatory, contractRequestIDstring)
						_, err := db.Query(query)
						if err != nil {
							fmt.Println("Unable to update database in relation to signatories")
						}
					}
				}
			}
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
		}
		tpl.ExecuteTemplate(res, "contractvalue.html", display)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
