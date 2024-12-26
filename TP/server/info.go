package server

import (
	"TP/structs"
	"TP/utils"
)

func InfoRequest(value []byte) []byte {
	var response []byte
	var request structs.InfoRequest

	request, message, err := request.Decode(value, games)
	if err != nil {
		print(err)
	}

	playerKey := players[request.PlayerId].PublicKey
	if utils.VerifySignature(&playerKey, message, request.Signature) {
		var info structs.InfoResponse

		if request.ValidMoves {
			info.Move = "Coups valides: \n" + GetValidMoves(games[request.GameId].Game)
		} else {
			info.Move, err = GetBestMove(games[request.GameId].Game)
			if err != nil {
				print(err)
			}
			info.Move = "Meilleur coup : " + info.Move
		}

		response, err = info.Encode(privateKey, games[request.GameId].EncryptionKey)
		if err != nil {
			print(err)
		}

	}

	return response
}
