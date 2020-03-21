package ken

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type GoogleTokenVerificationResponse struct {
	Error *string `json:"error"`
	Email string  `json:"email"`
	//Email string
}

func signPayload(payload string) string {
	key := []byte(os.Getenv("CRYPTKEY"))
	var h hash.Hash

	h = hmac.New(sha256.New, key)

	h.Write([]byte(payload))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//VerifyToken Google Sign in
func verifyGoogleToken(token string) (bool, string) {
	var verifyTokenURL = "https://oauth2.googleapis.com/tokeninfo?id_token=" + token

	res, err := http.Get(verifyTokenURL)

	if err != nil {
		return false, ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	//	fmt.Println(string(body))
	fmt.Println("string(body)")
	fmt.Println(signPayload(string(body)))

	if err != nil {
		fmt.Println("Error in ioutil")
		fmt.Println(err)
		return false, ""
	}

	var gResponse GoogleTokenVerificationResponse

	err = json.Unmarshal(body, &gResponse)

	if err != nil {
		fmt.Println("ERror in unmarshla")
		fmt.Println(err)
		return false, ""
	}

	// fmt.Println(gResponse)
	// fmt.Println("gResponse")

	if nil != gResponse.Error {
		return false, ""
	}

	// TODO handle this error
	newUserAdded := AddNewUser(gResponse.Email)

	fmt.Println(newUserAdded)
	return true, string(body)

}

//RouteVerifyIDToken ^)^
func RouteVerifyIDToken(c *gin.Context) {
	fmt.Println("c.Request.Body")
	verified, _ := verifyGoogleToken(c.Query("token"))
	if verified {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "gc6",
			Value:    "url.QueryEscape(userInfo)",
			MaxAge:   262800,
			Path:     "/",
			Domain:   "localhost",
			SameSite: 1,
			Secure:   false,
			HttpOnly: false,
		})
		c.SetCookie("gc1", "someName", 60*60*24, "/", "google.com", false, false)
		c.SetCookie("gc2", "someName", 60*60*24, "/", "localhost", false, false)
		c.SetCookie("gc3", "someName", 262800, "/", "localhost", false, false)
		c.SetCookie("gc4", "test", 3600, "/", "localhost", false, false)

		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "gc5",
			Value:    "12378",
			Expires:  time.Now().Add(time.Hour),
			HttpOnly: false,
			Secure:   false,
		})
		SendResponse(c, "OK", "Verified")
	} else {
		SendResponse(c, "Error", "Not verified")
	}

}
