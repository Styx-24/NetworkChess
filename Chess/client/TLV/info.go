package TLV

import (
	"TP/client/backEnd"
	"TP/structs"
	"TP/utils"
)

func InfoResponse(value []byte) []byte {
	var buffer []byte

	var response structs.InfoResponse

	response, message, err := response.Decode(value, EncryptionKey)
	if err != nil {
		print(err)
	}

	if utils.VerifySignature(&ServerPublicKey, message, response.Signature) {
		println(response.Move)

		buffer = backEnd.SelectMove(Player, GameId, EncryptionKey, IsPlayingSolo)
	}

	return buffer
}
