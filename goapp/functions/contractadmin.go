package functions

import (
	"fmt"
	"log"
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

type finalisedContract struct {
	ID       int
	Contract string
}

//ValueApproval to be exported
func ValueApproval(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	if myUser.Rights == "contractadmin" {
		results, err := Db.Query("SELECT Id, SigningEntity, CounterpartyName, Business, ContractValue FROM contracts_db.Contracts WHERE ContractValue IS NOT NULL AND FinanceTax = 'Pending' AND SeniorFinance IS NULL")
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
						_, err := Db.Query(query)
						if err != nil {
							fmt.Println("Unable to update database in relation to signatories")
						}
					}
				}
			}
			SendEmail("testtechnology.93@gmail.com")
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
		}
		Tpl.ExecuteTemplate(res, "contractvalue.html", display)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}

//ArchiveContract is to be exported
func ArchiveContract(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	if myUser.Rights == "contractadmin" {
		results, err := Db.Query("SELECT Id, Contract FROM contracts_db.Contracts WHERE Archived = 'Pending' AND Finalised = 'Signed'")
		if err != nil {
			log.Fatal(err)
		}
		display := []finalisedContract{}
		var contract finalisedContract
		for results.Next() {
			err := results.Scan(&contract.ID, &contract.Contract)
			if err != nil {
				fmt.Println("Unable to scan into contract variable")
			}
			display = append(display, contract)
		}

		if req.Method == http.MethodPost {
			contractRequestIDstring := req.FormValue("contractrequestid")
			contractRequestIDint, err := strconv.Atoi(contractRequestIDstring)
			if err != nil {
				fmt.Println("invalid, unapproved or no contract ID")
			}
			for _, v := range display {
				if v.ID == contractRequestIDint {
					query := fmt.Sprintf("UPDATE Contracts SET Archived = 'Yes' WHERE Id='%s'", contractRequestIDstring)
					_, err := Db.Query(query)
					if err != nil {
						fmt.Println("Unable to update database column 'Archived' to 'Yes'")
					}
				}
			}
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
		}
		Tpl.ExecuteTemplate(res, "archivecontract.html", display)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page.")
	}

}
