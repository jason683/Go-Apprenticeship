package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// parts of the function are taken from https://tutorialedge.net/golang/go-file-upload-tutorial/
func uploadFile(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(res, "draft.html", nil)

	if req.Method == http.MethodPost {
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
		fmt.Fprintf(res, "Successfully Uploaded File\n")
	}
}
