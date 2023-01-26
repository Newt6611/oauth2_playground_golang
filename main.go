package main

/* reference blog:
 * https://yushuanhsieh.github.io/post/2018-08-25-go-google-oauth/
 * https://medium.com/@pliutau/getting-started-with-oauth2-in-go-2c9fae55d187
**/

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleAuthCofig *oauth2.Config
)

func init() {
	googleAuthCofig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "YOUR_CLIENT_ID",
		ClientSecret: "YOUR_CLIENT_SECRET",
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	redirect_url := googleAuthCofig.AuthCodeURL("random_state")
	c.Redirect(http.StatusTemporaryRedirect, redirect_url)
}

func GoogleLoginCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "random_state" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("state unmatched"))
		return
	}

	// use "code" to get access token
	code := c.Query("code")

	// get access token
	token, err := googleAuthCofig.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	// print access token info
	c.JSON(http.StatusOK, token)

	//
	//
	//
	//
	// use access token to get userinfo
	client := googleAuthCofig.Client(context.Background(), token)
	res, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	defer res.Body.Close()
	rawData, _ := ioutil.ReadAll(res.Body)
	fmt.Println("userinfo: " + string(rawData))
	//
	// examle userinfo: {
	//    "sub": "123123124124312312",
	//    "name": "Newt6611",
	//    "given_name": "Newt6611",
	//    "picture": "https://lh3.googleusercontent.com/a/trihyoenrtiolg-c",
	//    "email": "guoching130@gmail.com",
	//    "email_verified": true,
	//    "locale": "zh-TW"
	// }
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/googleLogin", GoogleLogin)
	r.GET("/callback", GoogleLoginCallback)
	r.Run()
}
