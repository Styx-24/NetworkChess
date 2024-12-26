package client

import (
	"TP/structs"
	"TP/utils"
)

func GameResponse(value []byte) []byte {
	var buffer []byte

	var response structs.GameResponse

	response, message, err := response.Decode(value)

	if err != nil {
		print(err)
	} else {

		if utils.VerifySignature(&serverPublicKey, message, response.Signature) {
			println(response.Status)
			println("Équipe: " + utils.GetTeamString(response.Team))

			team = int(response.Team)
			encryptionKey = response.EncryptionKey
			gameId = response.GameId

			if team == response.TurnOf {

				buffer = SelectMove()

			}

		}
	}

	return buffer
}
