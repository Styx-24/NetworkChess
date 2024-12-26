package client

import (
	"TP/structs"
	"TP/utils"
)

func ActionResponse(value []byte) []byte {
	var buffer []byte

	var response structs.ActionResponse

	response, message, err := response.Decode(value, encryptionKey)

	if err != nil {
		print(err)
	} else {
		if utils.VerifySignature(&serverPublicKey, message, response.Signature) {
			println(response.Message)

			if response.MoveWasValid {
				println("Équipe: " + utils.GetTeamString(team))
			}

			if team == response.TurnOf {

				buffer = SelectMove()

			}

			if response.GameHasEnded {
				buffer, IsAPausedGame = GameSelection(player)
			}

		}
	}

	return buffer
}
