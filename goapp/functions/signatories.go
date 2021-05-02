package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type contract struct {
	ID            int
	Contract      string
	SeniorFinance string
}

//ShowContracts is to be exported
func ShowContracts(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	if myUser.Rights == "signatory" {
		//to add on new column (signed?)
		results, err := Db.Query("SELECT Id, Contract, SeniorFinance FROM contracts_db.Contracts WHERE Finalised = 'Pending' AND Contract IS NOT NULL")
		if err != nil {
			fmt.Println(err)
		}
		//display is a list of contracts
		display := []contract{}
		var row contract
		for results.Next() {
			err := results.Scan(&row.ID, &row.Contract, &row.SeniorFinance)
			if err != nil {
				fmt.Println(err)
			}
			display = append(display, row)
		}
		if req.Method == http.MethodPost {

			contractRequestIDstring := req.FormValue("contractrequestid")
			contractRequestIDint, err := strconv.Atoi(contractRequestIDstring)
			if err != nil {
				fmt.Println(err)
			}

			req.ParseMultipartForm(10 << 20)
			file, _, err := req.FormFile("myFile")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			tempFile, err := ioutil.TempFile("signed-documents", "upload-*.pdf")
			if err != nil {
				fmt.Println(err)
			}
			defer tempFile.Close()

			fileBytes, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
			}

			tempFile.Write(fileBytes)
			//fmt.Fprintf(res, "Successfully Uploaded File\n")

			filepath := tempFile.Name()

			for _, v := range display {
				if v.ID == contractRequestIDint {
					signed := "Yes"
					actionTime := time.Now().Format(time.RFC3339)
					_, err := Db.Query("UPDATE Contracts SET Contract = ?, Finalised = ?, ActionTime = ?, Archived = 'Pending' WHERE ID=?", filepath, signed, actionTime, contractRequestIDstring)
					if err != nil {
						fmt.Println(err)
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
		Tpl.ExecuteTemplate(res, "signatory.html", display)
	} else {
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
	}
}
