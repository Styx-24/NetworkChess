package client

import (
	"TP/structs"
	"TP/utils"
)

func GameComfirmationRequest(value []byte) []byte {
	var buffer []byte

	var request structs.GameComfirmationRequest

	request, message, err := request.Decode(value)
	if err != nil {
		println(err)
	} else {
		if utils.VerifySignature(&serverPublicKey, message, request.Signature) {
			println(request.Message)
			option := ComfirmationPromt()

			var response structs.GameComfirmationResponse
			response.Answer = option == 1
			response.GameId = gameId
			response.PlayerId = player.Id

			buffer, err = response.Encode(player.PrivateKey)
			if err != nil {
				println(err)
			}

		} else {
			println("signature invalide")
		}
	}

	return buffer
}

func GameComfirmationResponse(value []byte) []byte {
	var buffer []byte

	var response structs.GameComfirmationResponse

	response, message, err := response.Decode(value)
	if err != nil {
		println(err)
	} else {
		if utils.VerifySignature(&serverPublicKey, message, response.Signature) {
			if response.Answer {
				println("L'adversaire a accepter le match")
			} else {
				println("L'adversaire a refuser le match")
				buffer, IsAPausedGame = GameSelection(player)
			}

		}
	}

	return buffer
}
