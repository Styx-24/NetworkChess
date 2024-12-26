package server

import (
	"TP/structs"
	"database/sql"
	"log"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Bd structs.Bd

func DbCreation() error {
	exempleString := "CREATE TABLE IF NOT EXISTS Games(" +
		"ID  TEXT PRIMARY KEY," +
		"Player1 TEXT," +
		"Player2 TEXT," +
		"Fen TEXT," +
		"Turn INTEGER)"

	bd, err := sql.Open("sqlite3", "server.db")
	if err != nil {
		return err
	}

	if _, err := bd.Exec(exempleString); err != nil {
		return err
	}

	Bd = structs.Bd{Bd: bd}
	return nil
}

func InsertNewGame(game structs.Game) {

	queryString := "INSERT INTO Games(ID, Player1, Player2, Fen, Turn) VALUES (?,?,?,?,?)"

	db, err := sql.Open("sqlite3", "server.db")
	if err != nil {
		println(err.Error())
	}

	db.Exec(queryString, game.Id.String(), game.Player1.String(), game.Player2.String(), game.Game.FEN(), game.Turn)

}

func GetGames(playerId uuid.UUID) []structs.Game {

	db, err := sql.Open("sqlite3", "server.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queryString := "SELECT * FROM Games WHERE Player1 LIKE ? OR Player2 LIKE ?"

	rows, err := db.Query(queryString, playerId.String(), playerId.String())
	if err != nil {
		log.Fatal(err)
	}

	var games []structs.Game

	for rows.Next() {
		var gameBd structs.GameBd
		err = rows.Scan(&gameBd.Id, &gameBd.Player1, &gameBd.Player2, &gameBd.FEN, &gameBd.Turn)
		if err != nil {
			log.Fatal(err)
		}
		games = append(games, gameBd.DecodeGame())
	}

	rows.Close()
	return games
}

func GetGame(gameId uuid.UUID) structs.Game {
	db, err := sql.Open("sqlite3", "server.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queryString := "SELECT * FROM Games WHERE ID LIKE ?"

	rows, err := db.Query(queryString, gameId.String())
	if err != nil {
		log.Fatal(err)
	}

	var gameBd structs.GameBd

	if rows.Next() {
		err = rows.Scan(&gameBd.Id, &gameBd.Player1, &gameBd.Player2, &gameBd.FEN, &gameBd.Turn)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		println("no reslut for game in data base")
	}

	rows.Close()

	return gameBd.DecodeGame()
}

func DeleteGame(gameId uuid.UUID) {
	db, err := sql.Open("sqlite3", "server.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	queryString := "DELETE FROM Games WHERE ID LIKE ?"

	_, err = db.Exec(queryString, gameId.String())
	if err != nil {
		log.Fatal(err)
	}

}
