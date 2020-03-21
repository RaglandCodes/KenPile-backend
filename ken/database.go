package ken

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

//CheckIfUserExists returns true|false and an error
func CheckIfUserExists(emailToFind string) (bool, error) {
	db, err := sql.Open("postgres", os.Getenv("DBU"))

	if err != nil {
		return false, errors.New("Could't open DB")
	}
	fmt.Println(emailToFind + "is the email")

	rows, err := db.Query("SELECT COUNT(*) FROM users WHERE email = $1", emailToFind)

	defer rows.Close()

	if nil != err {
		fmt.Println("ERRor in query")
		return false, errors.New("Error in query")
	}

	for rows.Next() {
		var count int32
		if err := rows.Scan(&count); err != nil {
			fmt.Println(err)
			return false, errors.New("Could't scan")
		}

		if 0 == count {
			return false, nil
		} else if 1 == count {
			return true, nil
		} else {
			return false, errors.New("number of users with " + emailToFind + " is neither 0 nor 1")
		}

	}

	return false, errors.New("Some unexpected behavior in CheckIfUserExists")
}

//AddNewUser returns ERROR EXISTS or ADDED
func AddNewUser(email string) string {
	fmt.Println("in add new user")
	userExists, err := CheckIfUserExists(email)

	if nil != err {
		fmt.Println("error in userexists", err)
		return "ERROR"
	}

	if userExists {
		return "EXISTS"
	}

	db, err := sql.Open("postgres", os.Getenv("DBU"))

	insertStatment, err := db.Prepare("INSERT INTO users (email) VALUES ($1)")

	if err != nil {
		fmt.Println("error in prepare addnewuser", err)
		return "ERROR"
	}

	insertRes, err := insertStatment.Exec(email)

	if err != nil {
		fmt.Println("error in exec addnewuser", err)
		return "ERROR"
	}

	fmt.Println(insertRes)

	return "ADDED"
}

//RouteCreateNewNote returns OK or ERROR
func RouteCreateNewNote(c *gin.Context) {

	fmt.Println(c.Request.Body)
	fmt.Println(" bdy nn")
	fmt.Println(c.Cookie("gc1"))
	fmt.Println(c.Cookie("gc2"))
	fmt.Println(c.Cookie("gc3"))
	fmt.Println(c.Cookie("gc5"))
	fmt.Println(c.Cookie("gc6"))
	c4, err := c.Cookie("gc4")

	if nil != err {
		fmt.Println("c4 error")
		fmt.Println(err)
	}

	fmt.Println(c4)

	_, err = sql.Open("postgres", os.Getenv("DBU"))

	if nil != err {
		fmt.Println(err)
		fmt.Println("error in db open new note")
		SendResponse(c, "ERROR", "Error in opening database")
	}

	//BOMB
	SendResponse(c, "OK", "Created")

}

//Insert a random
func Insert() string {

	db, err := sql.Open("postgres", os.Getenv("DBU"))

	rows, err := db.Query("SELECT text FROM notes")
	if err != nil {
		return "ERROR"
	}
	defer rows.Close()
	res := ""
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			fmt.Println(err)
			return "ERROR"
		}
		res += email
		fmt.Println(email)
	}

	fmt.Println(err)
	return res
}
