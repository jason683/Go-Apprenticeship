package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//UploadFile : parts of the function was extracted from https://tutorialedge.net/golang/go-file-upload-tutorial/
func UploadFile(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}

	myUser := GetUser(res, req)
	if myUser.Rights == "legal" {
		results, err := Db.Query("SELECT Id, SigningEntity, CounterpartyName, Business, ContractType, ContractValue, Region, EffectiveDate, TerminationDate, BackgroundPurpose, CounterpartyContactInfo, Requester FROM contracts_db.Contracts WHERE ApproveStatus='Approve' AND Contract IS NULL")
		if err != nil {
			fmt.Println(err)
		}
		//display variable will contain a list of all approved contract requests
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

			contractRequestIDstring := req.FormValue("contractrequestid")
			contractRequestIDint, err := strconv.Atoi(contractRequestIDstring)
			if err != nil {
				fmt.Println("invalid, unapproved or no contract ID")
			}

			req.ParseMultipartForm(10 << 20)
			file, _, err := req.FormFile("myFile")
			if err != nil {
				fmt.Println("oh no")
				fmt.Println(err)
				return
			}
			defer file.Close()

			tempFile, err := ioutil.TempFile("temp-documents", "upload-*.pdf")
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}

			tempFile.Write(fileBytes)

			filepath := tempFile.Name()

			for _, v := range display {
				if v.ID == contractRequestIDint {
					actionTime := time.Now().Format(time.RFC3339)
					query := fmt.Sprintf("UPDATE Contracts SET Contract = '%s', ActionTime='%s' WHERE ID='%s'", filepath, actionTime, contractRequestIDstring)
					_, err := Db.Query(query)
					if err != nil {
						fmt.Println("Unable to update Contracts database")
					}
				}
			}
			emailAddress, err := Db.Query("SELECT Email FROM Users WHERE Rights = 'contractadmin'")
			if err != nil {
				fmt.Println(err)
			}
			var email string
			for emailAddress.Next() {
				err := emailAddress.Scan(&email)
				if err != nil {
					fmt.Println(err)
				}
				SendEmail(email)
			}
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
		}
		Tpl.ExecuteTemplate(res, "draft.html", display)
	} else {
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
	}

}
