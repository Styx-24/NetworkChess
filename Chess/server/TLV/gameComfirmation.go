package TLV

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
		playerKey := Players[comfirmation.PlayerId].PublicKey
		if utils.VerifySignature(&playerKey, message, comfirmation.Signature) {
			if comfirmation.Answer {

				var gameResponse structs.GameResponse

				delete(GameMatchMaking, comfirmation.PlayerId)
				gameResponse.GameId = comfirmation.GameId
				gameResponse.EncryptionKey = Games[comfirmation.GameId].EncryptionKey
				gameResponse.Status = Games[gameResponse.GameId].Game.Position().Board().Draw()
				gameResponse.Team = 1
				gameResponse.TurnOf = 1

				if Games[gameResponse.GameId].Player2Connexion != nil {
					gameResponse.Team = 2
					response, err = gameResponse.Encode(PrivateKey)
					if err != nil {
						println(err)
					}

					Games[gameResponse.GameId].Player2Connexion.Write(response)

					gameResponse.Team = 1

				}

				response, err = gameResponse.Encode(PrivateKey)
				if err != nil {
					println(err)
				}

			} else {
				comfirmationBuffer, err := comfirmation.Encode(PrivateKey)
				if err != nil {
					println(err)
				}
				if Games[comfirmation.GameId].Player2Connexion != nil {
					Games[comfirmation.GameId].Player2Connexion.Write(comfirmationBuffer)
				}
			}
		}
	}

	return response
}
