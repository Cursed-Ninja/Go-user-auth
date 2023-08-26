// GoogleOauth Related Functions are defined here

package googleoauth

import (
	"encoding/json"
	"net/http"
	"user-auth/internal/config/viper"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
)

type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Initialize the googleOauthConfig variable
func init() {
	googleOauthConfig.ClientID = viper.Get("google_oauth.client_id")
	googleOauthConfig.ClientSecret = viper.Get("google_oauth.client_secret")
	googleOauthConfig.RedirectURL = viper.Get("google_oauth.redirect_uri")
}

// Returns the url to redirect the user to google login page
func GoogleOauth() string {
	url := googleOauthConfig.AuthCodeURL("", oauth2.AccessTypeOffline)
	return url
}

// Helper function to get the user info from the google api
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

// Returns the user info from the google api
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
