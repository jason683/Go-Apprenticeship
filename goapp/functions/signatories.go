package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type contract struct {
	ID       int
	Contract string
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
		results, err := Db.Query("SELECT Id, Contract FROM contracts_db.Contracts WHERE SeniorFinance ='yes' AND Contract IS NOT NULL")
		if err != nil {
			fmt.Println("Cannot extract contract file")
		}
		//display is a list of contracts
		display := []contract{}
		var row contract
		for results.Next() {
			err := results.Scan(&row.ID, &row.Contract)
			if err != nil {
				fmt.Println("Cannot scan into row")
			}
			display = append(display, row)
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
					signed := "Signed"
					query := fmt.Sprintf("UPDATE Contracts SET Contract = '%s', Finalised = '%s' WHERE ID='%s'", filepath, signed, contractRequestIDstring)
					_, err := Db.Query(query)
					if err != nil {
						fmt.Println("Unable to update Contracts database")
					}
				}
			}
			SendEmail("testtechnology.93@gmail.com")
			http.Redirect(res, req, "/directory", http.StatusSeeOther)
		}
		Tpl.ExecuteTemplate(res, "signatory.html", display)
	} else {
		fmt.Fprintf(res, "You are not authorised to view this page")
	}
}
