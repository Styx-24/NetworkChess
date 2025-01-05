package TLV

import (
	"TP/client/backEnd"
	"TP/structs"
	"TP/utils"
)

func ActionResponse(value []byte) []byte {
	var buffer []byte

	var response structs.ActionResponse

	response, message, err := response.Decode(value, EncryptionKey)

	if err != nil {
		print(err)
	} else {
		if utils.VerifySignature(&ServerPublicKey, message, response.Signature) {
			println(response.Message)

			if response.MoveWasValid {
				println("Team: " + utils.GetTeamString(Team))
			}

			if Team == response.TurnOf {

				buffer = backEnd.SelectMove(Player, GameId, EncryptionKey, IsPlayingSolo)

			}

			if response.GameHasEnded {
				buffer, IsAPausedGame, IsPlayingSolo = backEnd.GameSelection(Player)
			}

		}
	}

	return buffer
}
