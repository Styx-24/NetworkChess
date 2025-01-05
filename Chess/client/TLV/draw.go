package TLV

import (
	"TP/client/backEnd"
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
		if utils.VerifySignature(&ServerPublicKey, message, request.Signature) {
			println(request.Message)
			option := backEnd.ComfirmationPromt()

			var response structs.DrawResponse
			response.Answer = option == 1
			response.GameId = GameId
			response.PlayerId = Player.Id

			buffer, err = response.Encode(Player.PrivateKey)
			if err != nil {
				println(err)
			}
		}
	}

	return buffer
}
