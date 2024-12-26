package client

import (
	"TP/structs"
	"TP/utils"
)

func InfoResponse(value []byte) []byte {
	var buffer []byte

	var response structs.InfoResponse

	response, message, err := response.Decode(value, encryptionKey)
	if err != nil {
		print(err)
	}

	if utils.VerifySignature(&serverPublicKey, message, response.Signature) {
		println(response.Move)

		buffer = SelectMove()
	}

	return buffer
}
