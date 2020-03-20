package ken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GoogleTokenVerificationResponse struct {
	Error *string `json:"error"`
	Email string  `json:"email"`
	//Email string
}

//VerifyToken Google Sign in
func verifyGoogleToken(token string) bool {
	var verifyTokenURL = "https://oauth2.googleapis.com/tokeninfo?id_token=" + token

	res, err := http.Get(verifyTokenURL)

	if err != nil {
		//SendResponse(c, 500, "Login failed!")
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	// fmt.Println(string(body))
	if err != nil {
		fmt.Println("ERror in ioutil")
		fmt.Println(err)
		return false
	}

	var gResponse GoogleTokenVerificationResponse

	err = json.Unmarshal(body, &gResponse)

	if err != nil {
		fmt.Println("ERror in unmarshla")
		fmt.Println(err)
		return false
	}

	// fmt.Println(gResponse)
	// fmt.Println("gResponse")

	if nil != gResponse.Error {
		return false
	}

	// TODO add user to DB
	newUserAdded := AddNewUser(gResponse.Email)

	fmt.Println(newUserAdded)
	return true

}

//RouteVerifyIDToken ^)^
func RouteVerifyIDToken(c *gin.Context) {
	fmt.Println("c.Request.Body")

	if verifyGoogleToken(c.Query("token")) {
		SendResponse(c, "OK", "Verified")
	} else {
		SendResponse(c, "Error", "Not verified")
	}

}
