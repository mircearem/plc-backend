package db

import (
	"database/sql"
	"errors"
	"fmt"
	auth "plc-backend/Auth"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseError struct {
	Message string `json:"message"`
}

func InsertUser(file string, user *auth.User) error {
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		return errors.New("Error opening database")
	}

	defer db.Close()

	// Format query string
	query := fmt.Sprintf(`
		INSERT INTO users 
		(username, password, email, admin)
		VALUES
		('%s', '%s', '%s', '%s');`,
		user.Username, user.Password, user.Email, user.Admin)

	insert, _ := db.Prepare(query)

	_, err = insert.Exec()

	if err != nil {
		return err
	}

	return nil
}

func FindUser(f string, u *auth.User) (*auth.User, error) {
	// Open the database
	db, err := sql.Open("sqlite3", f)

	if err != nil {
		return nil, errors.New("Error opening database")
	}

	defer db.Close()

	// Format the query string
	query := fmt.Sprintf(`SELECT * FROM users WHERE username='%s' OR email='%s' LIMIT 1;`, u.Username, u.Email)

	row, err := db.Query(query)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	// Parse the results of the query
	user := auth.User{}

	if row.Next() {
		if err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email, &user.Admin); err != nil {
			return nil, err
		}
	}

	// Check if result contains any data
	if user.Username == "" {
		return nil, nil
	}

	return &user, nil
}

func UpdateUser(f string, u *auth.User) error {
	db, err := sql.Open("sqlite3", f)

	if err != nil {
		return errors.New("Error opening database")
	}

	defer db.Close()

	return nil
}

func DeleteUser(file string, user *auth.User) error {
	return nil
}
