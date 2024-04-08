package controller

import (
	"Login-with-Outh/config"
	"context"
	"io/ioutil"
	"net/http"
)

func GoogleLogin(res http.ResponseWriter, req *http.Request) {
	googleConfig := config.SetupCon()
	url := googleConfig.AuthCodeURL("randomstate")

	http.Redirect(res, req, url, http.StatusSeeOther)
}

func GoogleCallback(res http.ResponseWriter, req *http.Request) {
	state := req.URL.Query().Get("state")
	if state != "randomstate" {
		http.Error(res, "State mismatch", http.StatusBadRequest)
		return
	}

	code := req.URL.Query().Get("code")
	if code == "" {
		http.Error(res, "Code not found in request", http.StatusBadRequest)
		return
	}

	googleConfig := config.SetupCon()

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(res, "Failed to exchange code for token", http.StatusInternalServerError)
		return
	}

	client := googleConfig.Client(context.Background(), token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(res, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}
	defer userInfoResp.Body.Close()

	if userInfoResp.StatusCode != http.StatusOK {
		http.Error(res, "Failed to fetch user data: "+userInfoResp.Status, userInfoResp.StatusCode)
		return
	}

	userData, err := ioutil.ReadAll(userInfoResp.Body)
	if err != nil {
		http.Error(res, "Failed to read user data", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(userData)
}
