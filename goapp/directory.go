package main

import (
	"net/http"
)

//directory

func directory(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := getUser(res, req)
	//concurrency issue on mapTypeRights
	if myUser.Rights == "bizrequester" {
		mapTypeRights["bizrequester"] = "hello requester"
		tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "bizrequester")
		return
	} else if myUser.Rights == "bizowner" {
		mapTypeRights["bizowner"] = "hello owner"
		tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "bizowner")
		return
	} else if myUser.Rights == "legal" {
		mapTypeRights["legal"] = "hello legal"
		tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "legal")
		return
	} else if myUser.Rights == "financetax" {
		mapTypeRights["financetax"] = "hello finance and tax"
		tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "financetax")
		return
	} else if myUser.Rights == "contractadmin" {
		mapTypeRights["contractadmin"] = "hello contract admin"
		tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "contractadmin")
		return
	} else {
		tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
	}
}
