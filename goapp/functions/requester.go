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
			newBusiness := req.FormValue("newbusiness")
			contractType := req.FormValue("contracttype")
			contractValue := req.FormValue("contractvalue")
			region := req.FormValue("region")
			effectiveDate := req.FormValue("effectivedate")
			terminationDate := req.FormValue("terminationdate")
			backgroundPurpose := req.FormValue("backgroundpurpose")
			counterpartyContactInfo := req.FormValue("counterpartycontactinfo")

			businessOwner := req.FormValue("businessowner")
			approveStatus := "Pending"

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
			_, err := strconv.Atoi(contractValue)
			if err != nil {
				errorMessage["input0"] = "Contract value has to be an integer"
				Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
				delete(errorMessage, "input0")
				return
			}

			timeAction := time.Now().Format(time.RFC3339)
			//NULL can be used to circumvent the int auto increment in sql
			Query := fmt.Sprintf("INSERT INTO Contracts (SigningEntity, CounterpartyName, Business, ContractType, ContractValue, Region, EffectiveDate, TerminationDate, BackgroundPurpose, CounterpartyContactInfo, Requester, BusinessOwner, ApproveStatus, ActionTime) VALUES ('%s', '%s', '%s', '%s', '%v', '%s', '%v', '%v', '%s', '%s', '%s', '%s', '%s', '%s')", signingEntity, counterpartyName, business, contractType, contractValue, region, effectiveDate, terminationDate, backgroundPurpose, counterpartyContactInfo, myUser.Username, businessOwner, approveStatus, timeAction)
			_, err = Db.Query(Query)
			if err != nil {
				fmt.Println(err)
			}

			//SendEmail("testtechnology.93@gmail.com")
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
			return
		}
		Tpl.ExecuteTemplate(res, "requestform.html", mapBusiness)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
