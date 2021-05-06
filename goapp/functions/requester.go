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

		if req.Method == http.MethodPost {
			errorMessage := map[string]string{}

			signingEntity := req.FormValue("signingentity")
			counterpartyName := req.FormValue("counterpartyname")
			business := req.FormValue("business")
			contractType := req.FormValue("contracttype")
			contractValue := req.FormValue("contractvalue")
			region := req.FormValue("region")
			effectiveDate := req.FormValue("effectivedate")
			terminationDate := req.FormValue("terminationdate")
			backgroundPurpose := req.FormValue("backgroundpurpose")
			counterpartyContactInfo := req.FormValue("counterpartycontactinfo")
			others := req.FormValue("others")

			businessOwner := req.FormValue("businessowner")
			approveStatus := "Pending"

			if signingEntity == "" || counterpartyName == "" || contractType == "" || businessOwner == "" || effectiveDate == "" || terminationDate == "" {
				errorMessage["input1"] = "Did you miss out entering any of the compulsory fields?"
				Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
				delete(errorMessage, "input1")
				return
			}
			if business == "" && others == "" {
				errorMessage["nobusiness"] = "You have not entered any value for the business type field"
				Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
				delete(errorMessage, "nobusiness")
				return
			}
			if contractValue != "" {
				_, err := strconv.Atoi(contractValue)
				if err != nil {
					errorMessage["input0"] = "Contract value has to be an integer"
					Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
					delete(errorMessage, "input0")
					return
				}
			}
			if business != "" && business != "Games & E-Commerce" && business != "Digital Finance" {
				errorMessage["input2"] = "You need to key in a valid business type"
				Tpl.ExecuteTemplate(res, "requestform.html", errorMessage)
				delete(errorMessage, "input2")
				return
			}
			if contractValue == "" {
				contractValue = "0"
			}
			if backgroundPurpose == "" {
				backgroundPurpose = "NA"
			}
			if region == "" {
				region = "NA"
			}
			if counterpartyContactInfo == "" {
				counterpartyContactInfo = "NA"
			}
			timeAction := time.Now().Format(time.RFC3339)
			//NULL can be used to circumvent the int auto increment in sql
			_, err := Db.Query("INSERT INTO Contracts (SigningEntity, CounterpartyName, Business, ContractType, ContractValue, Region, EffectiveDate, TerminationDate, BackgroundPurpose, CounterpartyContactInfo, Requester, BusinessOwner, ApproveStatus, ActionTime) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", signingEntity, counterpartyName, business, contractType, contractValue, region, effectiveDate, terminationDate, backgroundPurpose, counterpartyContactInfo, myUser.Username, businessOwner, approveStatus, timeAction)
			if err != nil {
				fmt.Println(err)
			}

			emailAddress, err := Db.Query("SELECT Email FROM Users WHERE Username = ?", businessOwner)
			if err != nil {
				fmt.Println(err)
			}
			var email string
			for emailAddress.Next() {
				err := emailAddress.Scan(&email)
				if err != nil {
					fmt.Println(err)
				}
			}
			SendEmail(email)
			relationMap[myUser.Username] = "Yes"
			http.Redirect(res, req, "/result", http.StatusSeeOther)
			return
		}
		Tpl.ExecuteTemplate(res, "requestform.html", mapBusiness)
	} else {
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
	}
}
