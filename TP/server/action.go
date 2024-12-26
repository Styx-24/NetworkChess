package server

import (
	"TP/structs"
	"TP/utils"

	"github.com/google/uuid"
)

func ActionRequest(value []byte) []byte {
	var response []byte
	var request structs.ActionRequest

	request, message, err := request.Decode(value, games)

	if err != nil {
		println(err)
	} else {

		playerKey := players[request.PlayerId].PublicKey

		if utils.VerifySignature(&playerKey, message, request.Signature) {
			var actionResponse structs.ActionResponse

			err = Move(games[request.GameId].Game, request.Move)
			if err != nil {
				print(err.Error())

				actionResponse.MoveWasValid = false
				actionResponse.GameHasEnded = false
				actionResponse.Message = "Coup Invalide"
				actionResponse.TurnOf = games[request.GameId].Turn

				response, err = actionResponse.Encode(privateKey, games[request.GameId].EncryptionKey)
				if err != nil {
					println(err)
				}

			} else {

				println("Coup Valid")
				actionResponse.MoveWasValid = true

				outcome, method := CheckVictory(games[request.GameId].Game)
				actionResponse.GameHasEnded = outcome != 0

				if !actionResponse.GameHasEnded && games[request.GameId].Player2 == uuid.Nil {
					games[request.GameId] = games[request.GameId].SetTurn(2)

					AIMove(games[request.GameId].Game)

					outcome, method = CheckVictory(games[request.GameId].Game)
					actionResponse.GameHasEnded = outcome != 0
				}

				if actionResponse.GameHasEnded {
					actionResponse.Message = utils.GetVictoryMessage(outcome) + " " + method
					actionResponse.TurnOf = 0
					actionResponse.MoveWasValid = false
				} else {
					actionResponse.Message = games[request.GameId].Game.Position().Board().Draw()

					if games[request.GameId].Turn == 1 {
						actionResponse.TurnOf = 2
						games[request.GameId] = games[request.GameId].SetTurn(2)
					} else {
						actionResponse.TurnOf = 1
						games[request.GameId] = games[request.GameId].SetTurn(1)
					}
				}

				response, err = actionResponse.Encode(privateKey, games[request.GameId].EncryptionKey)
				if err != nil {
					println(err)
				}

				if request.PlayerId == games[request.GameId].Player1 {
					if games[request.GameId].Player2Connexion != nil {
						games[request.GameId].Player2Connexion.Write(response)
					}
				} else {
					if games[request.GameId].Player2Connexion != nil {
						games[request.GameId].Player1Connexion.Write(response)
					}
				}

				if actionResponse.GameHasEnded {
					delete(games, request.GameId)
				}

			}

		}
	}

	return response
}
