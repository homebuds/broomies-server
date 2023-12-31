package service

import (
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
)

var google_config *oauth2.Config

type googleUser struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Gender        string `json:"gender"`
	Hd            string `json:"hd"`
}

// func getGoogleOauthURL() (*oauth2.Config, string) {
// 	options := CreateClientOptions("google", "https://ginoauth-example.herokuapp.com/callback/google")

// 	google_config = &oauth2.Config{
// 		ClientID:     options.getID(),
// 		ClientSecret: options.getSecret(),
// 		RedirectURL:  options.getRedirectURL(),
// 		Scopes: []string{
// 			"https://www.googleapis.com/auth/userinfo.email",
// 			"https://www.googleapis.com/auth/userinfo.profile",
// 		},
// 		Endpoint: google.Endpoint,
// 	}

// 	state := GenerateState()
// 	return google_config, state
// }

// func GoogleOauthLogin(ctx *gin.Context) {
// 	config, state := getGoogleOauthURL()
// 	redirectURL := config.AuthCodeURL(state)

// 	session := sessions.Default(ctx)
// 	session.Set("state", state)
// 	err := session.Save()
// 	if err != nil {
// 		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	ctx.Redirect(http.StatusSeeOther, redirectURL)
// }

// func GoogleCallBack(ctx *gin.Context) {
// 	session := sessions.Default(ctx)
// 	state := session.Get("state")
// 	if state != ctx.Query("state") {
// 		_ = ctx.AbortWithError(http.StatusUnauthorized, StateError)
// 		return
// 	}

// 	code := ctx.Query("code")
// 	token, err := google_config.Exchange(ctx, code)
// 	if err != nil {
// 		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
// 		return
// 	}

// 	client := google_config.Client(context.TODO(), token)
// 	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
// 	if err != nil {
// 		_ = ctx.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}
// 	defer userInfo.Body.Close()

// 	info, err := ioutil.ReadAll(userInfo.Body)
// 	if err != nil {
// 		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	var user googleUser
// 	err = json.Unmarshal(info, &user)
// 	if err != nil {
// 		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	redirectURL, err := url.Parse(IsLoginURL)
// 	if err != nil {
// 		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	query, err := url.ParseQuery(redirectURL.RawQuery)
// 	if err != nil {
// 		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
// 		return
// 	}

// 	query.Add("email", user.Email)
// 	query.Add("name", user.Name)
// 	query.Add("source", "google")
// 	redirectURL.RawQuery = query.Encode()

// 	ctx.Redirect(http.StatusSeeOther, redirectURL.String())
// }
