package client

import (
	"TP/structs"
	"TP/utils"

	"github.com/google/uuid"
)

func MatchMakingResponse(value []byte) []byte {
	var buffer []byte

	var matchMakingResponse structs.MatchMakingResponse

	matchMakingResponse, message, err := matchMakingResponse.Decode(value)

	if err != nil {
		print(err)
	} else {
		if utils.VerifySignature(&serverPublicKey, message, matchMakingResponse.Signature) {
			id := OpponentSelection(matchMakingResponse)

			if id == 0 {
				buffer, IsAPausedGame = GameSelection(player)
			} else {

				var GameRequest structs.GameRequest

				if IsAPausedGame {
					GameRequest.GameId = matchMakingResponse.IDs[id-1]
				} else {
					GameRequest.PlayerId = matchMakingResponse.IDs[id-1]
					GameRequest.GameId = uuid.Nil
				}

				GameRequest.OponentId = player.Id

				buffer, err = GameRequest.Encode(player.PrivateKey)
				if err != nil {
					println(err)
				}
			}
		}
	}

	return buffer
}
