package server

import (
	"TP/structs"
	"TP/utils"

	"github.com/notnil/chess"
)

func DrawRequest(value []byte) []byte {
	var request structs.DrawRequest
	var response []byte

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

func DrawResponse(value []byte) []byte {

	var queryResponse structs.DrawResponse

	var response []byte

	queryResponse, message, err := queryResponse.Decode(value)
	if err != nil {
		println(err)
	} else {
		playerKey := players[queryResponse.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, queryResponse.Signature) {

			if queryResponse.Answer {
				games[queryResponse.GameId].Game.Draw(chess.DrawOffer)

				var actionResponse structs.ActionResponse

				actionResponse.GameHasEnded = true
				actionResponse.Message = "Partie null"
				actionResponse.MoveWasValid = false
				actionResponse.TurnOf = 0

				response, err = actionResponse.Encode(privateKey, games[queryResponse.GameId].EncryptionKey)
				if err != nil {
					println(err)
				}

				if games[queryResponse.GameId].Player1 == queryResponse.PlayerId {
					games[queryResponse.GameId].Player2Connexion.Write(response)
				} else {
					games[queryResponse.GameId].Player1Connexion.Write(response)
				}

			} else {
				println("Denied")
			}

		}
	}

	return response

}
