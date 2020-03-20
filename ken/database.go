package ken

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
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
		//res += email
		fmt.Println(count)

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
