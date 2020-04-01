package ken

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type GoogleTokenVerificationResponse struct {
	Error *string `json:"error"`
	Email string  `json:"email"`
	//Email string
}

// GetEmailFromToken verifies token and returns email
func GetEmailFromToken(token string) (string, error) {

	tokenArray := strings.Split(token, ".")

	payload, decodeError := base64.StdEncoding.DecodeString(tokenArray[0])

	if nil != decodeError {
		fmt.Println("decodeError")
		fmt.Println(decodeError)
		return "", errors.New("Decode error")
	}

	signedPayload := signPayload(string(payload))

	if signedPayload == tokenArray[1] {

		var p GoogleTokenVerificationResponse

		unmarshalError := json.Unmarshal(payload, &p)

		if unmarshalError != nil {
			return "", errors.New("Unmarshal Error")
		}

		return p.Email, nil

	}

	return "", errors.New("Couldn't verify")
}
func signPayload(payload string) string {
	// TODO take []byte as input instead of string
	key := []byte(os.Getenv("CRYPTKEY"))
	var h hash.Hash

	h = hmac.New(sha256.New, key)

	h.Write([]byte(payload))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

//VerifyToken Google Sign and if successful, also return auth token
func verifyGoogleToken(token string) (bool, string) {
	verifyTokenURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + token

	res, err := http.Get(verifyTokenURL)

	if err != nil {
		return false, ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	// TODO rename this body variable name.

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

	if nil != gResponse.Error {
		return false, ""
	}

	// TODO handle this error
	newUserAdded := AddNewUser(gResponse.Email)

	fmt.Println(newUserAdded)
	return true, base64.StdEncoding.EncodeToString(body) + "." + signPayload(string(body))

}

//RouteVerifyIDToken verified token and sets cookie
func RouteVerifyIDToken(c *gin.Context) {
	fmt.Println("c.Request.Body")
	verified, authToken := verifyGoogleToken(c.Query("token"))
	if verified {
		c.SetCookie("aatt", authToken, 100000, "/", "localhost", false, true)
		fmt.Println(authToken)
		SendResponse(c, "OK", "Verified")
	} else {
		SendResponse(c, "Error", "Not verified")
	}

}
