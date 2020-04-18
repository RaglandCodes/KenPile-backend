package ken

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type updateBody struct {
	NoteID   string `json:"id"`
	NoteText string `json:"text"`
}

//RouteUpdateNote ^)^
func RouteUpdateNote(c *gin.Context) {
	if VerifyCookieToken(c) == false {
		SendResponse(c, "ERROR", "Auth Error")
		return
	}

	sub, subErr := GetSubFromCookie(c)

	if nil != subErr {
		fmt.Println(subErr)
		fmt.Println("subErr ^ ")
		SendResponse(c, "ERROR", "Auth Error")
		return
	}

	db, err := sql.Open("postgres", os.Getenv("DBU"))
	defer db.Close()

	if nil != err {
		fmt.Println(err)
		fmt.Println("error in db open new note")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	updateStatment, err := db.Prepare("UPDATE notes SET text = $1 WHERE owner = $2 AND id = $3")

	if nil != err {
		fmt.Println(err)
		fmt.Println("error in db open update note")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	ub := updateBody{}

	if bindErr := c.ShouldBindBodyWith(&ub, binding.JSON); bindErr != nil {
		fmt.Println(bindErr)
		fmt.Println("error in bind")
		SendResponse(c, "ERROR", "Bind error")
		return
	}

	_, updateError := updateStatment.Exec(ub.NoteText, sub, ub.NoteID)

	if nil != updateError {
		fmt.Println(updateError)
		fmt.Println("error in db open new note")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	SendResponse(c, "OK", "Updated")

}
