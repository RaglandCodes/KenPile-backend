package ken

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type updateBody struct {
	NoteID    string `json:"id"`
	NoteText  string `json:"text"`
	NoteDelta string `json:"delta"`
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

	updateStatment, err := db.Prepare("UPDATE notes SET text = $1, delta = $4, last_updated_at = NOW() WHERE owner = $2 AND id = $3")

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

	_, updateError := updateStatment.Exec(ub.NoteText, sub, ub.NoteID, ub.NoteDelta)

	if nil != updateError {
		fmt.Println(updateError)
		fmt.Println("error in db open new note")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	SendResponse(c, "OK", "Updated")

}

// RouteFetchAllNotes will return all notes of that user
func RouteFetchAllNotes(c *gin.Context) {
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

	db, dbErr := sql.Open("postgres", os.Getenv("DBU"))
	defer db.Close()

	if nil != dbErr {
		fmt.Println(dbErr)
		fmt.Println("error in db open new note")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	//SELECT * FROM notes WHERE owner = $1

	rows, rowsErr := db.Query("SELECT id, delta FROM notes WHERE owner = $1", sub)
	defer rows.Close()

	if nil != rowsErr {
		fmt.Println(rowsErr)
		fmt.Println("rowsErr ^ ")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	type NoteRow struct {
		ID    string
		Delta sql.NullString
	}

	var results []NoteRow

	for rows.Next() {

		noteRow := NoteRow{}

		if scanErr := rows.Scan(&noteRow.ID, &noteRow.Delta); scanErr != nil {
			log.Fatal(scanErr)
		}

		results = append(results, noteRow)
	}

	m, mErr := json.Marshal(results)

	if nil != mErr {
		fmt.Println(mErr)
		fmt.Println("mErr ^ ")
		SendResponse(c, "ERROR", "Marhsal Error")
	}

	SendResponse(c, "OK", string(m))
}

//RouteFetchNote returns the Delta of one note
func RouteFetchNote(c *gin.Context) {
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

	db, dbErr := sql.Open("postgres", os.Getenv("DBU"))
	defer db.Close()

	if nil != dbErr {
		fmt.Println(dbErr)
		fmt.Println("error in db open new note")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	var delta sql.NullString

	row := db.QueryRow("SELECT delta FROM notes WHERE id=$1 AND owner=$2", c.Query("id"), sub)
	scanErr := row.Scan(&delta)

	if nil != scanErr {
		fmt.Println(scanErr)
		fmt.Println("scanErr ^ ")
		SendResponse(c, "ERROR", "DB Error")
		return
	}

	SendResponse(c, "OK", delta.String)

}
