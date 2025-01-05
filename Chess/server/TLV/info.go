package TLV

import (
	"TP/server/backEnd"
	"TP/structs"
	"TP/utils"
)

func InfoRequest(value []byte) []byte {
	var response []byte
	var request structs.InfoRequest

	request, message, err := request.Decode(value, Games)
	if err != nil {
		print(err)
	}

	playerKey := Players[request.PlayerId].PublicKey
	if utils.VerifySignature(&playerKey, message, request.Signature) {
		var info structs.InfoResponse

		if request.ValidMoves {
			info.Move = Games[request.GameId].Game.Position().Board().Draw() + "\n Valid moves: \n" + backEnd.GetValidMoves(Games[request.GameId].Game)
		} else {
			info.Move, err = backEnd.GetBestMove(Games[request.GameId].Game)
			if err != nil {
				print(err)
			}
			info.Move = Games[request.GameId].Game.Position().Board().Draw() + "\n Best move : " + info.Move
		}

		response, err = info.Encode(PrivateKey, Games[request.GameId].EncryptionKey)
		if err != nil {
			print(err)
		}

	}

	return response
}
