package db

import (
	"database/sql"
	"fmt"
	util "plc-backend/Utils"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func FindUserForSignup(file string, user *util.User) (bool, *util.DatabaseError) {
	db, err := sql.Open("sqlite3", file)
	defer db.Close()

	// Error opening the database file
	if err != nil {
		err := util.DatabaseError{
			Message: "Unable to open database file",
		}
		return false, &err
	}

	// Format the query string
	query := fmt.Sprintf(`SELECT username,email FROM users WHERE username='%s' OR email='%s';`, user.Username, user.Email)

	// Execute query
	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		err := util.DatabaseError{
			Message: "Unable to execute query",
		}
		return false, &err
	}

	// Check if there are any rows returned by the query
	cols, _ := rows.Columns()

	if len(cols) == 0 {
		return false, nil
	}

	return true, nil
}

func FindUserForLogin(file string, username string) (util.LoginResponse, *util.DatabaseError) {
	result := util.LoginResponse{}

	// Error opening the database file
	db, err := sql.Open("sqlite3", file)
	defer db.Close()

	// Error opening the database
	if err != nil {
		err := util.DatabaseError{
			Message: "Unable to open database file",
		}
		return result, &err
	}

	// Format the query string
	query := fmt.Sprintf(`SELECT username,password FROM users WHERE username='%s';`, username)

	rows, _ := db.Query(query)
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&result.Username, &result.Password); err != nil {
			err := util.DatabaseError{
				Message: "Unable to parse query results",
			}
			return result, &err
		}
	}

	return result, nil
}

func InsertUser(file string, user *util.User) error {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return err
	}

	// Format query string
	query := fmt.Sprintf(`
		INSERT INTO users 
		(username, password, email, admin)
		VALUES
		('%s', '%s', '%s', '%s');`,
		user.Username, user.Password, user.Email, strconv.FormatBool(*user.Admin))

	insert, _ := db.Prepare(query)

	_, err = insert.Exec()

	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func UpdateUser(file string, user *util.User, column string, param string) error {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return err
	}

	// Format query string
	query := fmt.Sprintf(`
    UPDATE users
    SET %s = '%s' WHERE
    username = '%s';`, column, param, user.Username)

	// Execute query
	update, _ := db.Prepare(query)

	_, err = update.Exec()

	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func DeleteUser(file string, user *util.User) error {
	return nil
}
