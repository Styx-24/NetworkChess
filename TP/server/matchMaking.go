package server

import (
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

	playerKey := players[matchMakingRequest.ID].PublicKey
	if utils.VerifySignature(&playerKey, message, matchMakingRequest.Signature) {
		var matchMakingResponse structs.MatchMakingResponse

		matchMakingResponse.IDs = make([]uuid.UUID, 0)
		matchMakingResponse.Names = make([]string, 0)

		if matchMakingRequest.IsAPausedGame {

			games := GetGames(matchMakingRequest.ID)
			for _, game := range games {
				matchMakingResponse.IDs = append(matchMakingResponse.IDs, game.Id)
				if players[game.Player2].Name == "" {
					matchMakingResponse.Names = append(matchMakingResponse.Names, players[game.Player1].Name+" "+players[game.Player1].LastName+" vs IA")
				} else {
					matchMakingResponse.Names = append(matchMakingResponse.Names, players[game.Player1].Name+" "+players[game.Player1].LastName+" vs "+players[game.Player2].Name+" "+players[game.Player2].LastName)
				}
			}

		} else {

			for _, value := range gameMatchMaking {
				matchMakingResponse.IDs = append(matchMakingResponse.IDs, value.Player1)
				matchMakingResponse.Names = append(matchMakingResponse.Names, players[value.Player1].Name+" "+players[value.Player1].LastName)
			}
		}

		response, err = matchMakingResponse.Encode(privateKey)
		if err != nil {
			println(err)
		}

	}

	return response
}
