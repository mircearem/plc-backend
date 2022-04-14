package db

import (
	"database/sql"
	"fmt"
	util "plc-backend/Utils"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func FindUserByStruct(file string, user *util.User) (*util.User, error) {
	result := new(util.User)

	db, err := sql.Open("sqlite3", file)

	// Error opening the database file
	if err != nil {
		return nil, err
	}

	// Format the query string
	query := fmt.Sprintf(`
		SELECT username,email 
		FROM users 
		WHERE username='%s' OR email='%s';
		`, user.Username, user.Email)

	// Execute query
	rows, _ := db.Query(query)

	// Check if there are any rows returned by the query
	if !rows.Next() {
		return nil, nil
	}

	// Interpret query
	for rows.Next() {
		err := rows.Scan(&result)

		if err != nil {
			return nil, err
		}
	}

	rows.Close()
	db.Close()

	return result, nil
}

func FindUserByUsername(file string, username string) (*util.User, error) {
	result := new(util.User)

	// Error opening the database file
	db, err := sql.Open("sqlite3", file)

	// Error
	if err != nil {
		return nil, err
	}

	// Format the query string
	query := fmt.Sprintf(`SELECT username FROM users WHERE username='%s';`, username)

	rows, _ := db.Query(query)

	// Check if there are any rows returned by the query
	if !rows.Next() {
		return nil, nil
	}

	// Interpret query
	for rows.Next() {
		err := rows.Scan(&result)

		if err != nil {
			return nil, err
		}
	}

	rows.Close()
	db.Close()

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
