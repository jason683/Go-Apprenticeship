package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// parts of the function taken from https://tutorialedge.net/golang/go-file-upload-tutorial/
func uploadFile(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}

	results, err := db.Query("SELECT Id, SigningEntity, CounterpartyName, Business, Requester FROM contracts_db.Contracts WHERE ApproveStatus='Approve'")
	if err != nil {
		fmt.Println("Something has happened")
	}
	//display variable will contain a list of all approved contract requests
	display := []contractRequest{}
	var request contractRequest
	for results.Next() {
		err := results.Scan(&request.ID, &request.SigningEntity, &request.CounterpartyName, &request.Business, &request.Requester)
		if err != nil {
			fmt.Println("Something has happened in request variable")
		}
		display = append(display, request)
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
		//fmt.Fprintf(res, "Successfully Uploaded File\n")

		filepath := tempFile.Name()

		for _, v := range display {
			if v.ID == contractRequestIDint {
				query := fmt.Sprintf("UPDATE Contracts SET Contract = '%s' WHERE ID='%s'", filepath, contractRequestIDstring)
				_, err := db.Query(query)
				if err != nil {
					fmt.Println("Unable to update Contracts database")
				}
			}
		}
		http.Redirect(res, req, "/directory", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(res, "draft.html", display)
}
