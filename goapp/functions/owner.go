package functions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type contractRequest struct {
	ID                      int
	SigningEntity           string
	CounterpartyName        string
	Business                string
	ContractType            string
	ContractValue           int
	Region                  string
	EffectiveDate           string
	TerminationDate         string
	BackgroundPurpose       string
	CounterpartyContactInfo string
	Requester               string
}

type testing struct {
	EffectiveDate string
}

//ReviewRequest is to be exported
func ReviewRequest(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	if myUser.Rights == "bizowner" {
		results, err := Db.Query("SELECT Id, SigningEntity, CounterpartyName, Business, ContractType, ContractValue, Region, EffectiveDate, TerminationDate, BackgroundPurpose, CounterpartyContactInfo, Requester FROM contracts_db.Contracts WHERE BusinessOwner = ? AND ApproveStatus ='Pending'", myUser.Username)
		if err != nil {
			panic(err.Error())
		}
		//display variable will contain a list of all the pending contract requests
		display := []contractRequest{}
		var reviewRequest contractRequest
		for results.Next() {
			err := results.Scan(&reviewRequest.ID, &reviewRequest.SigningEntity, &reviewRequest.CounterpartyName, &reviewRequest.Business, &reviewRequest.ContractType, &reviewRequest.ContractValue, &reviewRequest.Region, &reviewRequest.EffectiveDate, &reviewRequest.TerminationDate, &reviewRequest.BackgroundPurpose, &reviewRequest.CounterpartyContactInfo, &reviewRequest.Requester)
			if err != nil {
				panic(err.Error())
			}
			display = append(display, reviewRequest)
		}
		if req.Method == http.MethodPost {
			if !AlreadyLoggedIn(req) {
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
						if v.ContractValue > 0 {
							actionTime := time.Now().Format(time.RFC3339)
							query := fmt.Sprintf("UPDATE Contracts SET ApproveStatus='%s', FinanceTax='Pending', ActionTime='%s' WHERE Id='%s'", contractRequestStatus, actionTime, contractRequestIDstring)
							_, err := Db.Query(query)
							if err != nil {
								panic(err.Error())
							}
						} else {
							actionTime := time.Now().Format(time.RFC3339)
							query := fmt.Sprintf("UPDATE Contracts SET ApproveStatus='%s', ActionTime='%s' WHERE Id='%s'", contractRequestStatus, actionTime, contractRequestIDstring)
							_, err := Db.Query(query)
							if err != nil {
								panic(err.Error())
							}
						}
					}
				}
			}
			SendEmail("testtechnology.93@gmail.com")
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
		}
		Tpl.ExecuteTemplate(res, "revrequest.html", display)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
