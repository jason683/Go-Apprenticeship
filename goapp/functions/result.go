package functions

import "net/http"

//Result allows for the display of an outcome page to be reflected after a user has completed an action.
func Result(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	mapResult := make(map[string]string)
	if myUser.Rights == "bizrequester" {
		if relationMap[myUser.Username] == "Yes" {
			mapResult["bizrequester"] = "You have successfully submitted the contract request."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "bizrequester")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	if myUser.Rights == "bizowner" {
		if relationMap[myUser.Username] == "Yes" {
			mapResult["bizowner"] = "You have successfully reviewed the contract request."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "bizowner")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	if myUser.Rights == "financetax" {
		if relationMap[myUser.Username] == "Yes" {
			mapResult["financetax"] = "You have successfully reviewed the contract request."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "financetax")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	if myUser.Rights == "legal" {
		if relationMap[myUser.Username] == "Yes" {
			mapResult["legal"] = "You have successfully uploaded the contract document."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "legal")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	if myUser.Rights == "contractadmin" {
		if relationMap[myUser.Username] == "Yes0" {
			mapResult["contractadmin"] = "You have successfully reviewed the contract request amount."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "contractadmin")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	if myUser.Rights == "contractadmin" {
		if relationMap[myUser.Username] == "Yes1" {
			mapResult["contractadmin"] = "The database now records the contract document as archived."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "contractadmin")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	if myUser.Rights == "contractadmin" {
		if relationMap[myUser.Username] == "Yes2" {
			mapResult["contractadmin"] = "You have sent an email to notify the action owner."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "contractadmin")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	if myUser.Rights == "signatory" {
		if relationMap[myUser.Username] == "Yes" {
			mapResult["signatory"] = "You have successfully uploaded the signed contract."
			Tpl.ExecuteTemplate(res, "result.html", mapResult)
			delete(mapResult, "signatory")
			relationMap[myUser.Username] = "No"
			return
		}
	}
	Tpl.ExecuteTemplate(res, "result.html", mapResult)
}
