package structs

import (
	"crypto/rsa"
	"database/sql"
	"net"

	"github.com/google/uuid"
	"github.com/notnil/chess"
)

type Bd struct {
	Bd *sql.DB
}

type User struct {
	Name       string
	LastName   string
	Id         uuid.UUID
	PrivateKey rsa.PrivateKey
	PublicKey  rsa.PublicKey
}

type Game struct {
	Id               uuid.UUID
	Player1          uuid.UUID
	Player2          uuid.UUID
	Player1Connexion net.Conn
	Player2Connexion net.Conn
	EncryptionKey    []byte
	FEN              string
	Game             *chess.Game
	Turn             int
}

func (h Game) SetTurn(turn int) Game {
	var tmp Game

	tmp.Turn = turn
	tmp.Id = h.Id
	tmp.Player1 = h.Player1
	tmp.Player2 = h.Player2
	tmp.Player1Connexion = h.Player1Connexion
	tmp.Player2Connexion = h.Player2Connexion
	tmp.EncryptionKey = h.EncryptionKey
	tmp.Game = h.Game
	tmp.FEN = h.FEN

	return tmp
}

type GameBd struct {
	Id      string
	Player1 string
	Player2 string
	FEN     string
	Turn    int
}

func (h GameBd) DecodeGame() Game {
	var tmp Game

	tmp.Turn = h.Turn
	tmp.Id, _ = uuid.Parse(h.Id)
	tmp.Player1, _ = uuid.Parse(h.Player1)
	tmp.Player2, _ = uuid.Parse(h.Player2)
	tmp.FEN = h.FEN

	return tmp
}
