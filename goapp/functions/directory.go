package functions

import (
	"net/http"
)

//Directory to be exported
func Directory(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	//concurrency issue on mapTypeRights
	if myUser.Rights == "bizrequester" {
		mapTypeRights["bizrequester"] = "hello requester"
		Tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "bizrequester")
		return
	} else if myUser.Rights == "bizowner" {
		mapTypeRights["bizowner"] = "hello owner"
		Tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "bizowner")
		return
	} else if myUser.Rights == "legal" {
		mapTypeRights["legal"] = "hello legal"
		Tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "legal")
		return
	} else if myUser.Rights == "financetax" {
		mapTypeRights["financetax"] = "hello finance and tax"
		Tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "financetax")
		return
	} else if myUser.Rights == "contractadmin" {
		mapTypeRights["contractadmin"] = "hello contract admin"
		Tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "contractadmin")
		return
	} else if myUser.Rights == "signatory" {
		mapTypeRights["signatory"] = "hello signatory"
		Tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
		delete(mapTypeRights, "signatory")
		return
	} else {
		Tpl.ExecuteTemplate(res, "directory.html", mapTypeRights)
	}
}
