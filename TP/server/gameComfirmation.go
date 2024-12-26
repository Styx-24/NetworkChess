package server

import (
	"TP/structs"
	"TP/utils"
)

func GameComfirmationResponse(value []byte) []byte {

	var comfirmation structs.GameComfirmationResponse
	var response []byte

	comfirmation, message, err := comfirmation.Decode(value)
	if err != nil {
		println(err)
	} else {
		playerKey := players[comfirmation.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, comfirmation.Signature) {
			if comfirmation.Answer {

				var gameResponse structs.GameResponse

				delete(gameMatchMaking, comfirmation.PlayerId)
				gameResponse.GameId = comfirmation.GameId
				gameResponse.EncryptionKey = games[comfirmation.GameId].EncryptionKey
				gameResponse.Status = games[gameResponse.GameId].Game.Position().Board().Draw()
				gameResponse.Team = 1
				gameResponse.TurnOf = 1

				if games[gameResponse.GameId].Player2Connexion != nil {
					gameResponse.Team = 2
					response, err = gameResponse.Encode(privateKey)
					if err != nil {
						println(err)
					}

					games[gameResponse.GameId].Player2Connexion.Write(response)

					gameResponse.Team = 1

				}

				response, err = gameResponse.Encode(privateKey)
				if err != nil {
					println(err)
				}

			} else {
				comfirmationBuffer, err := comfirmation.Encode(privateKey)
				if err != nil {
					println(err)
				}
				if games[comfirmation.GameId].Player2Connexion != nil {
					games[comfirmation.GameId].Player2Connexion.Write(comfirmationBuffer)
				}
			}
		}
	}

	return response
}
