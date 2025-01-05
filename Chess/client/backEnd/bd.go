package backEnd

import (
	"TP/structs"
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Bd structs.Bd

func DbCreation() error {
	exempleString := "CREATE TABLE IF NOT EXISTS Users(" +
		"ID  TEXT PRIMARY KEY," +
		"Name TEXT," +
		"LastName TEXT)"

	bd, err := sql.Open("sqlite3", "client.db")
	if err != nil {
		return err
	}

	if _, err := bd.Exec(exempleString); err != nil {
		return err
	}

	Bd = structs.Bd{Bd: bd}
	return nil
}

func InsertNewUser(user structs.User) {

	queryString := "INSERT INTO Users(ID, Name, LastName) VALUES (?,?,?)"

	db, err := sql.Open("sqlite3", "client.db")
	if err != nil {

	}

	db.Exec(queryString, user.Id.String(), user.Name, user.LastName)

}

func GetUsers() []structs.User {

	db, err := sql.Open("sqlite3", "./client.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM Users")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var users []structs.User

	for rows.Next() {
		var user structs.User
		var idString string
		err := rows.Scan(&idString, &user.Name, &user.LastName)
		if err != nil {
			log.Fatal(err)
		}
		user.Id, err = uuid.Parse(idString)
		users = append(users, user)
	}

	return users

}
