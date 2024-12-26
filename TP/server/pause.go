package server

import (
	"TP/structs"
	"TP/utils"

	"github.com/google/uuid"
)

func PauseRequest(value []byte) []byte {

	var request structs.PauseRequest
	var response []byte = make([]byte, 0)

	request, message, err := request.Decode(value)
	if err != nil {
		println(err)
	} else {

		playerKey := players[request.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, request.Signature) {

			requestBuffer, err := request.Encode(privateKey)
			if err != nil {
				println(err)
			}

			if games[request.GameId].Player1 == request.PlayerId {
				if games[request.GameId].Player2Connexion != nil {
					games[request.GameId].Player2Connexion.Write(requestBuffer)
				} else {
					response = pauseGame(request.GameId, request.PlayerId)
				}

			} else {
				if games[request.GameId].Player1Connexion != nil {
					games[request.GameId].Player1Connexion.Write(requestBuffer)
				}
			}
		}
	}

	return response

}

func PauseResponse(value []byte) []byte {

	var queryResponse structs.PauseResponse
	var response []byte

	queryResponse, message, err := queryResponse.Decode(value)
	if err != nil {
		println(err)
	} else {
		playerKey := players[queryResponse.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, queryResponse.Signature) {

			if queryResponse.Answer {
				response = pauseGame(queryResponse.GameId, queryResponse.PlayerId)
			} else {
				println("Denied")
			}

		}
	}

	return response

}

func pauseGame(gameId uuid.UUID, playerId uuid.UUID) []byte {
	var response []byte
	var err error
	var actionResponse structs.ActionResponse

	actionResponse.GameHasEnded = true
	actionResponse.Message = "Partie en pause"
	actionResponse.MoveWasValid = false
	actionResponse.TurnOf = 0

	response, err = actionResponse.Encode(privateKey, games[gameId].EncryptionKey)
	if err != nil {
		println(err)
	}

	if games[gameId].Player1 == playerId {
		if games[gameId].Player2Connexion != nil {
			games[gameId].Player2Connexion.Write(response)
		}
	} else {
		if games[gameId].Player1Connexion != nil {
			games[gameId].Player1Connexion.Write(response)
		}
	}

	InsertNewGame(games[gameId])

	delete(games, gameId)

	return response
}
