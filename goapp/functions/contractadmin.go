package functions

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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

type timing struct {
	ID     int
	Timing time.Time
}

type row struct {
	ID            int
	BusinessOwner string
	ApproveStatus string
	Contract      string
	FinanceTax    string
	Finalised     string
	Archived      string
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
						query := fmt.Sprintf("UPDATE Contracts SET SeniorFinance='%s', Finalised='Pending' WHERE Id='%s'", signatory, contractRequestIDstring)
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

//IdentifyOutdatedRequest is to be exported
func IdentifyOutdatedRequest(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	if myUser.Rights == "contractadmin" {
		results, err := Db.Query("SELECT ID, ActionTime FROM Contracts WHERE Archived IS NOT NULL AND Archived != 'Yes'")
		if err != nil {
			fmt.Println(err)
		}
		display := []timing{}
		var lastDone timing
		for results.Next() {
			err := results.Scan(&lastDone.ID, &lastDone.Timing)
			if err != nil {
				fmt.Println(err)
			}
			display = append(display, lastDone)
		}

		for _, v := range display {
			//if more than 7 days
			if time.Now().Sub(v.Timing).Hours() >= 168 {
				_, err := Db.Query(fmt.Sprintf("UPDATE Contracts SET Outdated = 'Yes' WHERE ID = '%v'", v.ID))
				if err != nil {
					fmt.Println(err)
				}
				//Tpl.ExecuteTemplate(res, "reminder.html", display)
			} else {
				_, err := Db.Query(fmt.Sprintf("UPDATE Contracts SET Outdated = 'No' WHERE ID = '%v'", v.ID))
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		displaySecond := []row{}
		var contract row
		dbquery, err := Db.Query("SELECT Id, BusinessOwner, ApproveStatus, Contract, FinanceTax, Finalised, Archived FROM Contracts WHERE Outdated = 'Yes'")
		if err != nil {
			fmt.Println(err)
		}
		for dbquery.Next() {
			err := dbquery.Scan(&contract.ID, &contract.BusinessOwner, &contract.ApproveStatus, &contract.Contract, &contract.FinanceTax, &contract.Finalised, &contract.Archived)
			if err != nil {
				fmt.Println(err)
			}
			displaySecond = append(displaySecond, contract)
		}
		Tpl.ExecuteTemplate(res, "outdatedcontract.html", displaySecond)

	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}

//EmailList is to be exported
func EmailList(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myUser := GetUser(res, req)
	if myUser.Rights == "contractadmin" {
		dbquery, err := Db.Query("SELECT Email FROM Users")
		if err != nil {
			fmt.Println(err)
		}
		var email string
		emailList := []string{}
		for dbquery.Next() {
			err := dbquery.Scan(&email)
			if err != nil {
				fmt.Println(err)
			}
			emailList = append(emailList, email)
		}

		if req.Method == http.MethodPost {
			email := req.FormValue("email")
			for _, v := range emailList {
				if v == email {
					SendEmail(v)
					http.Redirect(res, req, "/directory", http.StatusSeeOther)
				}
			}
			fmt.Fprintf(res, "No valid email was found.")
		}

		Tpl.ExecuteTemplate(res, "emaillist.html", emailList)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
