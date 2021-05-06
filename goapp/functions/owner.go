package functions

import (
	"fmt"
	"net/http"
	"strconv"
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
			fmt.Println(err)
		}
		//display variable will contain a list of all the pending contract requests
		display := []contractRequest{}
		var reviewRequest contractRequest
		for results.Next() {
			err := results.Scan(&reviewRequest.ID, &reviewRequest.SigningEntity, &reviewRequest.CounterpartyName, &reviewRequest.Business, &reviewRequest.ContractType, &reviewRequest.ContractValue, &reviewRequest.Region, &reviewRequest.EffectiveDate, &reviewRequest.TerminationDate, &reviewRequest.BackgroundPurpose, &reviewRequest.CounterpartyContactInfo, &reviewRequest.Requester)
			if err != nil {
				fmt.Println(err)
			}
			reviewRequest.EffectiveDate = reviewRequest.EffectiveDate[:10]
			reviewRequest.TerminationDate = reviewRequest.TerminationDate[:10]
			display = append(display, reviewRequest)
		}
		if req.Method == http.MethodPost {
			if !AlreadyLoggedIn(req) {
				http.Redirect(res, req, "/", http.StatusSeeOther)
			}
			contractRequestIDstring := req.FormValue("contractrequestid")
			contractRequestIDint, err := strconv.Atoi(contractRequestIDstring)
			if err != nil {
				fmt.Println(err)
			}
			//section below will change pending status to either approve or reject status
			contractRequestStatus := req.FormValue("approvereject")
			for _, v := range display {
				if v.ID == contractRequestIDint {
					if contractRequestStatus == "Approve" || contractRequestStatus == "Reject" {
						if v.ContractValue > 0 {
							actionTime := time.Now().Format(time.RFC3339)
							_, err := Db.Query("UPDATE Contracts SET ApproveStatus=?, FinanceTax='Pending', ActionTime=? WHERE Id=?", contractRequestStatus, actionTime, contractRequestIDstring)
							if err != nil {
								fmt.Println(err)
							}
							// emailAddress, err := Db.Query("SELECT Email FROM Users WHERE Rights = 'financetax'")
							// if err != nil {
							// 	fmt.Println(err)
							// }
							// var email string
							// for emailAddress.Next() {
							// 	err := emailAddress.Scan(&email)
							// 	if err != nil {
							// 		fmt.Println(err)
							// 	}
							// 	SendEmail(email)
							// }
						} else {
							actionTime := time.Now().Format(time.RFC3339)
							_, err := Db.Query("UPDATE Contracts SET ApproveStatus=?, FinanceTax='NA', ActionTime=? WHERE Id=?", contractRequestStatus, actionTime, contractRequestIDstring)
							if err != nil {
								fmt.Println(err)
							}
							// emailAddress, err := Db.Query("SELECT Email FROM Users WHERE Rights = 'legal'")
							// if err != nil {
							// 	fmt.Println(err)
							// }
							// var email string
							// for emailAddress.Next() {
							// 	err := emailAddress.Scan(&email)
							// 	if err != nil {
							// 		fmt.Println(err)
							// 	}
							// 	SendEmail(email)
							// }
						}
					}
				}
			}
			relationMap[myUser.Username] = "Yes"
			http.Redirect(res, req, "/result", http.StatusSeeOther)
		}
		Tpl.ExecuteTemplate(res, "revrequest.html", display)
	} else {
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
	}
}
