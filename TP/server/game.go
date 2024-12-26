package server

import (
	"TP/structs"
	"TP/utils"
	"crypto/rand"
	"net"

	"github.com/google/uuid"
)

func GameRequest(value []byte, c net.Conn) []byte {
	var response []byte

	var GameRequest structs.GameRequest

	GameRequest, message, err := GameRequest.Decode(value)

	if err != nil {
		println(err)
	} else {

		playerKey := players[GameRequest.PlayerId].PublicKey
		opponentKey := players[GameRequest.OponentId].PublicKey

		if utils.VerifySignature(&playerKey, message, GameRequest.Signature) || utils.VerifySignature(&opponentKey, message, GameRequest.Signature) {

			if GameRequest.GameId != uuid.Nil {
				response = PausedGameRequest(value, GameRequest, c)
			} else if GameRequest.OponentId == uuid.Nil {
				response = SoloGameRequest(value, GameRequest, c)
			} else {
				response = PVPGameRequest(value, GameRequest, c)
			}

		} else {
			println("Signature invalide")
		}

	}

	return response
}

func PausedGameRequest(value []byte, GameRequest structs.GameRequest, c net.Conn) []byte {
	var response []byte
	var err error
	var gameResponse structs.GameResponse
	gameResponse.GameId = uuid.New()

	if games[GameRequest.GameId].Player1 == uuid.Nil {

		game := GetGame(GameRequest.GameId)
		game.Game = LoadGame(game.FEN)

		var key = make([]byte, 16)

		_, err := rand.Read(key)
		if err != nil {
			println(err)
		}

		game.EncryptionKey = key

		if game.Player2 == uuid.Nil {

			gameResponse.EncryptionKey = key
			gameResponse.Team = 1
			gameResponse.GameId = GameRequest.GameId
			gameResponse.TurnOf = game.Turn
			gameResponse.Status = game.Game.Position().Board().Draw()

			DeleteGame(games[GameRequest.GameId].Id)

		} else {

			gameResponse.EncryptionKey = key
			gameResponse.Team = 0
			gameResponse.TurnOf = 3
			gameResponse.Status = "Attente d'un adversère"

		}

		if GameRequest.PlayerId == game.Player1 {
			game.Player1Connexion = c
		} else {
			game.Player2Connexion = c
		}

		games[GameRequest.GameId] = game

	} else {

		DeleteGame(games[GameRequest.GameId].Id)

		game := games[GameRequest.GameId]

		gameResponse.GameId = GameRequest.GameId
		gameResponse.EncryptionKey = game.EncryptionKey
		gameResponse.Status = game.Game.Position().Board().Draw()
		gameResponse.TurnOf = game.Turn

		if GameRequest.PlayerId == game.Player1 {
			game.Player1Connexion = c
		} else {
			game.Player2Connexion = c
		}

		games[GameRequest.GameId] = game

		if GameRequest.PlayerId == games[gameResponse.GameId].Player1 {
			gameResponse.Team = 2
			response, err = gameResponse.Encode(privateKey)
			if err != nil {
				println(err)
			}

			games[gameResponse.GameId].Player2Connexion.Write(response)

			gameResponse.Team = 1

		} else {
			gameResponse.Team = 1
			response, err = gameResponse.Encode(privateKey)
			if err != nil {
				println(err)
			}

			games[gameResponse.GameId].Player1Connexion.Write(response)

			gameResponse.Team = 2
		}

	}

	response, err = gameResponse.Encode(privateKey)
	if err != nil {
		println(err)
	}

	return response
}

func SoloGameRequest(value []byte, GameRequest structs.GameRequest, c net.Conn) []byte {
	var response []byte
	var gameResponse structs.GameResponse
	gameResponse.GameId = uuid.New()

	var key = make([]byte, 16)

	_, err := rand.Read(key)
	if err != nil {
	}

	gameResponse.EncryptionKey = key
	gameResponse.TurnOf = 1

	games[gameResponse.GameId] = structs.Game{Player1: players[GameRequest.PlayerId].Id, Player1Connexion: c, Player2: uuid.Nil, Id: gameResponse.GameId, EncryptionKey: key, Turn: 1, Game: GenerateGame()}

	gameResponse.Status = games[gameResponse.GameId].Game.Position().Board().Draw()
	gameResponse.Team = 1

	response, err = gameResponse.Encode(privateKey)
	if err != nil {
		println(err)
	}

	return response
}

func PVPGameRequest(value []byte, GameRequest structs.GameRequest, c net.Conn) []byte {
	var response []byte
	var gameResponse structs.GameResponse
	var err error

	gameResponse.GameId = uuid.New()

	if GameRequest.PlayerId == GameRequest.OponentId {

		var key = make([]byte, 16)

		_, err := rand.Read(key)
		if err != nil {
			println(err)
		}
		gameResponse.EncryptionKey = key
		gameResponse.Team = 0
		gameResponse.Status = "Attente d'un adversère"
		gameResponse.TurnOf = 2
		gameMatchMaking[GameRequest.PlayerId] = structs.Game{Player1: players[GameRequest.PlayerId].Id, Player1Connexion: c, Id: gameResponse.GameId, EncryptionKey: key}

	} else {

		games[gameMatchMaking[GameRequest.PlayerId].Id] = structs.Game{Player1: players[GameRequest.PlayerId].Id, Player1Connexion: gameMatchMaking[GameRequest.PlayerId].Player1Connexion, Player2: GameRequest.OponentId, Player2Connexion: c, Id: gameMatchMaking[GameRequest.PlayerId].Id, EncryptionKey: gameMatchMaking[GameRequest.PlayerId].EncryptionKey, Turn: 1, Game: GenerateGame()}
		gameResponse.GameId = gameMatchMaking[GameRequest.PlayerId].Id
		gameResponse.Status = "Attente de comfirmation avec l'adverère"
		gameResponse.EncryptionKey = gameMatchMaking[GameRequest.PlayerId].EncryptionKey
		gameResponse.Team = 2
		gameResponse.TurnOf = 1

		var comfirmation structs.GameComfirmationRequest
		comfirmation.Message = "Joueur " + players[GameRequest.PlayerId].Name + " " + players[GameRequest.PlayerId].LastName + " veux jouer contre vous, voulez-vous commencer le match?"
		comfirmationBuffer, err := comfirmation.Encode(privateKey)
		if err != nil {
			println(err)
		}

		games[gameResponse.GameId].Player1Connexion.Write(comfirmationBuffer)

	}

	response, err = gameResponse.Encode(privateKey)
	if err != nil {
		println(err)
	}

	return response
}
