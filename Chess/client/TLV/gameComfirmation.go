package TLV

import (
	"TP/client/backEnd"
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
		if utils.VerifySignature(&ServerPublicKey, message, request.Signature) {
			println(request.Message)
			option := backEnd.ComfirmationPromt()

			var response structs.GameComfirmationResponse
			response.Answer = option == 1
			response.GameId = GameId
			response.PlayerId = Player.Id

			buffer, err = response.Encode(Player.PrivateKey)
			if err != nil {
				println(err)
			}

		} else {
			println("Invalid Signature")
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
		if utils.VerifySignature(&ServerPublicKey, message, response.Signature) {
			if response.Answer {
				println("The opponent accepted the match")
			} else {
				println("The opponent refused the match")
				buffer, IsAPausedGame, IsPlayingSolo = backEnd.GameSelection(Player)
			}

		}
	}

	return buffer
}
