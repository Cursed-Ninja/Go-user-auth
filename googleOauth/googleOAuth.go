package googleoauth

import (
	"encoding/json"
	"net/http"
	"user-auth/config/viper"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
)

type UserInfo struct {
	Id    string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func init() {
	googleOauthConfig.ClientID = viper.Get("google_oauth.client_id")
	googleOauthConfig.ClientSecret = viper.Get("google_oauth.client_secret")
	googleOauthConfig.RedirectURL = viper.Get("google_oauth.redirect_uri")
}

func GoogleOauth() string {
	url := googleOauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
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

func Callback(r *http.Request) (*UserInfo, error) {
	code := r.URL.Query().Get("code")

	token, err := googleOauthConfig.Exchange(r.Context(), code)

	if err != nil {
		return nil, err
	}

	client := googleOauthConfig.Client(r.Context(), token)

	userInfo, err := getUserInfo(client)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}
