package TLV

import (
	"TP/server/backEnd"
	"TP/structs"
	"TP/utils"

	"github.com/google/uuid"
)

func ActionRequest(value []byte) []byte {
	var response []byte
	var request structs.ActionRequest

	request, message, err := request.Decode(value, Games)

	if err != nil {
		println(err)
	} else {

		playerKey := Players[request.PlayerId].PublicKey

		if utils.VerifySignature(&playerKey, message, request.Signature) {
			var actionResponse structs.ActionResponse

			err = backEnd.Move(Games[request.GameId].Game, request.Move)
			if err != nil {
				print(err.Error())

				actionResponse.MoveWasValid = false
				actionResponse.GameHasEnded = false
				actionResponse.Message = "Move invalid"
				actionResponse.TurnOf = Games[request.GameId].Turn

				response, err = actionResponse.Encode(PrivateKey, Games[request.GameId].EncryptionKey)
				if err != nil {
					println(err)
				}

			} else {

				actionResponse.MoveWasValid = true

				outcome, method := backEnd.CheckVictory(Games[request.GameId].Game)
				actionResponse.GameHasEnded = outcome != 0

				if !actionResponse.GameHasEnded && Games[request.GameId].Player2 == uuid.Nil {
					Games[request.GameId] = Games[request.GameId].SetTurn(2)

					backEnd.AIMove(Games[request.GameId].Game)

					outcome, method = backEnd.CheckVictory(Games[request.GameId].Game)
					actionResponse.GameHasEnded = outcome != 0
				}

				if actionResponse.GameHasEnded {
					actionResponse.Message = utils.GetVictoryMessage(outcome) + " " + method
					actionResponse.TurnOf = 0
					actionResponse.MoveWasValid = false
				} else {
					actionResponse.Message = Games[request.GameId].Game.Position().Board().Draw()

					if Games[request.GameId].Turn == 1 {
						actionResponse.TurnOf = 2
						Games[request.GameId] = Games[request.GameId].SetTurn(2)
					} else {
						actionResponse.TurnOf = 1
						Games[request.GameId] = Games[request.GameId].SetTurn(1)
					}
				}

				response, err = actionResponse.Encode(PrivateKey, Games[request.GameId].EncryptionKey)
				if err != nil {
					println(err)
				}

				if request.PlayerId == Games[request.GameId].Player1 {
					if Games[request.GameId].Player2Connexion != nil {
						Games[request.GameId].Player2Connexion.Write(response)
					}
				} else {
					if Games[request.GameId].Player2Connexion != nil {
						Games[request.GameId].Player1Connexion.Write(response)
					}
				}

				if actionResponse.GameHasEnded {
					delete(Games, request.GameId)
				}

			}

		}
	}

	return response
}
