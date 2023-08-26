package googleoauth

import (
	"encoding/json"
	"net/http"
	"user-auth/config/viper"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfigLogin = &oauth2.Config{
		ClientID:     viper.Get("google_oauth.common.client_id"),
		ClientSecret: viper.Get("google_oauth.common.client_secret"),
		RedirectURL:  viper.Get("google_oauth.login.redirect_uri"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	googleOauthConfigRegister = &oauth2.Config{
		ClientID:     viper.Get("google_oauth.common.client_id"),
		ClientSecret: viper.Get("google_oauth.common.client_secret"),
		RedirectURL:  viper.Get("google_oauth.register.redirect_uri"),
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
)

type UserInfo struct {
	Id    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
	// Add other fields as needed
}

func GoogleOauth(method string) string {
	var url string
	if method == "login" {
		url = googleOauthConfigLogin.AuthCodeURL("", oauth2.AccessTypeOffline)
	} else {
		url = googleOauthConfigRegister.AuthCodeURL("", oauth2.AccessTypeOffline)
	}
	return url
}

func getUserInfo(client *http.Client) (*UserInfo, error) {
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func Callback(r *http.Request, method string) (*UserInfo, error) {
	code := r.URL.Query().Get("code")

	var token *oauth2.Token
	var err error
	if method == "login" {
		token, err = googleOauthConfigLogin.Exchange(r.Context(), code)
	} else {
		token, err = googleOauthConfigRegister.Exchange(r.Context(), code)
	}

	if err != nil {
		return nil, err
	}

	var client *http.Client

	if method == "login" {
		client = googleOauthConfigLogin.Client(r.Context(), token)
	} else {
		client = googleOauthConfigRegister.Client(r.Context(), token)
	}

	userInfo, err := getUserInfo(client)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
