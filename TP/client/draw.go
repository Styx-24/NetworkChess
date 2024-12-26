package client

import (
	"TP/structs"
	"TP/utils"
)

func DrawRequest(value []byte) []byte {
	var buffer []byte

	var request structs.DrawRequest

	request, message, err := request.Decode(value)
	if err != nil {
		println(err)
	} else {
		if utils.VerifySignature(&serverPublicKey, message, request.Signature) {
			println(request.Message)
			option := ComfirmationPromt()

			var response structs.DrawResponse
			response.Answer = option == 1
			response.GameId = gameId
			response.PlayerId = player.Id

			buffer, err = response.Encode(player.PrivateKey)
			if err != nil {
				println(err)
			}
		}
	}

	return buffer
}
