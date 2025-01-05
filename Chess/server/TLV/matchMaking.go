package TLV

import (
	"TP/server/backEnd"
	"TP/structs"
	"TP/utils"

	"github.com/google/uuid"
)

func MatchMakingRequest(value []byte) []byte {
	var response []byte

	var matchMakingRequest structs.MatchMakingRequest

	matchMakingRequest, message, err := matchMakingRequest.Decode(value)
	if err != nil {
		println(err)
	}

	playerKey := Players[matchMakingRequest.ID].PublicKey
	if utils.VerifySignature(&playerKey, message, matchMakingRequest.Signature) {
		var matchMakingResponse structs.MatchMakingResponse

		matchMakingResponse.IDs = make([]uuid.UUID, 0)
		matchMakingResponse.Names = make([]string, 0)

		if matchMakingRequest.IsAPausedGame {

			Games := backEnd.GetGames(matchMakingRequest.ID)
			for _, game := range Games {
				matchMakingResponse.IDs = append(matchMakingResponse.IDs, game.Id)
				if Players[game.Player2].Name == "" {
					matchMakingResponse.Names = append(matchMakingResponse.Names, Players[game.Player1].Name+" "+Players[game.Player1].LastName+" vs IA")
				} else {
					matchMakingResponse.Names = append(matchMakingResponse.Names, Players[game.Player1].Name+" "+Players[game.Player1].LastName+" vs "+Players[game.Player2].Name+" "+Players[game.Player2].LastName)
				}
			}

		} else {

			for _, value := range GameMatchMaking {
				matchMakingResponse.IDs = append(matchMakingResponse.IDs, value.Player1)
				matchMakingResponse.Names = append(matchMakingResponse.Names, Players[value.Player1].Name+" "+Players[value.Player1].LastName)
			}
		}

		response, err = matchMakingResponse.Encode(PrivateKey)
		if err != nil {
			println(err)
		}

	}

	return response
}
