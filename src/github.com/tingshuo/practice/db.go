package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"fmt"
)

func logAndReturnError(err error) error{
	log.Println(err.Error())
	return err
}

func createUser(u *user, db *sql.DB) error{
	if err := validateInput(u, db); err != nil{
		return err
	}
	statement := "INSERT INTO users(name, age) VALUES(?, ?)"
	stmt, err := db.Prepare(statement)
	defer stmt.Close()
	if err != nil{
		logAndReturnError(err)
	}
	if _, err = stmt.Exec(u.Name, u.Age); err != nil{
		log.Println(fmt.Sprintf("User name: %s", u.Name))
		log.Println(fmt.Sprintf("User Age: %d", u.Age))
		logAndReturnError(err)
	}
	if err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID); err != nil {
		logAndReturnError(err)
	}
	return err
}

func deleteUser(u *user, db *sql.DB) error{
	if err := validateInput(u, db); err != nil{
		return err
	}
	statement := "DELETE FROM users WHERE id=?"
	stmt, err := db.Prepare(statement)
	defer stmt.Close()
	if err != nil{
		logAndReturnError(err)
	}
	if _, err = stmt.Exec(u.ID) ;err != nil{
		log.Println(fmt.Sprintf("User ID: %d", u.ID))
		logAndReturnError(err)
	}
	return err
}

func updateUser(u *user, db *sql.DB) error{
	if err := validateInput(u, db); err != nil{
		logAndReturnError(err)
	}
	statement := "UPDATE users SET name=?, age=? WHERE id=?"
	stmt, err := db.Prepare(statement)
	defer stmt.Close()
	if err != nil{
		logAndReturnError(err)
	}
	if _, err = stmt.Exec(u.Name, u.Age, u.ID) ;err != nil {
		log.Println(fmt.Sprintf("User ID: %d", u.ID))
		log.Println(fmt.Sprintf("User Name: %s", u.Name))
		log.Println(fmt.Sprintf("User Age: %d", u.Age))
		logAndReturnError(err)
	}
	return err
}

func getUser(u *user, db *sql.DB) error {
	if err := validateInput(u, db); err != nil{
		logAndReturnError(err)
	}
	statement := "SELECT name, ifnull(age,0) FROM users WHERE id=?"
	stmt, err := db.Prepare(statement)
	defer stmt.Close()
	if err != nil{
		log.Println(fmt.Sprintf("User ID: %d", u.ID))
		logAndReturnError(err)
	}
	if  err = stmt.QueryRow(u.ID).Scan(&u.Name, &u.Age); err != nil{
		log.Println(fmt.Sprintf("User ID: %d", u.ID))
		logAndReturnError(err)
	}
	return err
}

func getUsers(limit int, offset int, db *sql.DB) ([]user, error){
	if db == nil{
		msg := "db connection is nil"
		log.Print(msg)
		return nil, argError{msg}
	}
	if offset < 0 || limit < 0{
		msg := "user count and offset must be positive integers"
		log.Print(msg)
		return nil, argError{msg}
	}
	statement := "SELECT id, name, ifnull(age, 0) FROM users LIMIT ? OFFSET ?"
	stmt, err := db.Prepare(statement)
	defer stmt.Close()
	if err != nil{
		log.Print("Failed to prepare sql statement:" + statement)
		return nil, err
	}
	rows, err := stmt.Query(limit, offset)
	defer rows.Close()
	if err != nil{
		log.Print(fmt.Sprintf("Failed to query users with argument count:%d offset:%d", limit, offset))
		return nil, err
	}
	users := []user{}
	for rows.Next() {
		var u user
		if err := rows.Scan(&u.ID, &u.Name, &u.Age); err != nil {
			log.Print(fmt.Sprintf("Failed to query users with argument count:%d offset:%d", limit, offset))
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func validateInput(u *user, db *sql.DB) error{
	if u == nil{
		msg := "user is nil"
		log.Print(msg)
		return argError{msg}
	}
	if db == nil{
		msg := "db connection is nil"
		log.Print(msg)
		return argError{msg}
	}
	return nil
}

