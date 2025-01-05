package TLV

import (
	"TP/server/backEnd"
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

		playerKey := Players[request.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, request.Signature) {

			requestBuffer, err := request.Encode(PrivateKey)
			if err != nil {
				println(err)
			}

			if Games[request.GameId].Player1 == request.PlayerId {
				if Games[request.GameId].Player2Connexion != nil {
					Games[request.GameId].Player2Connexion.Write(requestBuffer)
				} else {
					response = pauseGame(request.GameId, request.PlayerId)
				}

			} else {
				if Games[request.GameId].Player1Connexion != nil {
					Games[request.GameId].Player1Connexion.Write(requestBuffer)
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
		playerKey := Players[queryResponse.PlayerId].PublicKey
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
	actionResponse.Message = "Game paused"
	actionResponse.MoveWasValid = false
	actionResponse.TurnOf = 0

	response, err = actionResponse.Encode(PrivateKey, Games[gameId].EncryptionKey)
	if err != nil {
		println(err)
	}

	if Games[gameId].Player1 == playerId {
		if Games[gameId].Player2Connexion != nil {
			Games[gameId].Player2Connexion.Write(response)
		}
	} else {
		if Games[gameId].Player1Connexion != nil {
			Games[gameId].Player1Connexion.Write(response)
		}
	}

	backEnd.InsertNewGame(Games[gameId])

	delete(Games, gameId)

	return response
}
