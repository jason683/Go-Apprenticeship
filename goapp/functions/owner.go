package functions

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type contractRequest struct {
	ID               int
	SigningEntity    string
	CounterpartyName string
	Business         string
	Requester        string
}

//ReviewRequest is to be exported
func ReviewRequest(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	results, err := Db.Query("SELECT Id, SigningEntity, CounterpartyName, Business, Requester FROM contracts_db.Contracts WHERE BusinessOwner = ? AND ApproveStatus ='Pending'", myUser.Username)
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
					query := fmt.Sprintf("UPDATE Contracts SET ApproveStatus='%s' WHERE Id='%s'", contractRequestStatus, contractRequestIDstring)
					_, err := Db.Query(query)
					if err != nil {
						panic(err.Error())
					}
				}
			}
		}
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
	}
	Tpl.ExecuteTemplate(res, "revrequest.html", display)
}
