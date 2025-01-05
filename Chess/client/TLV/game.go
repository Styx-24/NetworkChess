package TLV

import (
	"TP/client/backEnd"
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

		if utils.VerifySignature(&ServerPublicKey, message, response.Signature) {
			println(response.Status)
			println("Team: " + utils.GetTeamString(response.Team))

			Team = int(response.Team)
			EncryptionKey = response.EncryptionKey
			GameId = response.GameId

			if Team == response.TurnOf {

				buffer = backEnd.SelectMove(Player, GameId, EncryptionKey, IsPlayingSolo)

			}

		}
	}

	return buffer
}
