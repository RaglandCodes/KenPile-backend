package ken

import (
	"database/sql"
	"fmt"
	"os"
)

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
