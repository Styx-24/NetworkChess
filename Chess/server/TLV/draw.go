package TLV

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
		playerKey := Players[request.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, request.Signature) {

			requestBuffer, err := request.Encode(PrivateKey)
			if err != nil {
				println(err)
			}

			if Games[request.GameId].Player1 == request.PlayerId {
				if Games[request.GameId].Player2Connexion != nil {
					Games[request.GameId].Player2Connexion.Write(requestBuffer)
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

func DrawResponse(value []byte) []byte {

	var queryResponse structs.DrawResponse

	var response []byte

	queryResponse, message, err := queryResponse.Decode(value)
	if err != nil {
		println(err)
	} else {
		playerKey := Players[queryResponse.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, queryResponse.Signature) {

			if queryResponse.Answer {
				Games[queryResponse.GameId].Game.Draw(chess.DrawOffer)

				var actionResponse structs.ActionResponse

				actionResponse.GameHasEnded = true
				actionResponse.Message = "Null game"
				actionResponse.MoveWasValid = false
				actionResponse.TurnOf = 0

				response, err = actionResponse.Encode(PrivateKey, Games[queryResponse.GameId].EncryptionKey)
				if err != nil {
					println(err)
				}

				if Games[queryResponse.GameId].Player1 == queryResponse.PlayerId {
					Games[queryResponse.GameId].Player2Connexion.Write(response)
				} else {
					Games[queryResponse.GameId].Player1Connexion.Write(response)
				}

			} else {
				println("Denied")
			}

		}
	}

	return response

}
