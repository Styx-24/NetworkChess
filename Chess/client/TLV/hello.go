package TLV

import (
	"TP/client/backEnd"
	"TP/structs"
	"crypto/rsa"
)

func HelloResponse(value []byte) ([]byte, rsa.PublicKey) {
	var buffer []byte

	var helloResponse structs.HelloResponse
	ServerPublicKey, err := helloResponse.Decode(value)
	ServerPublicKey.Size()
	if err != nil {
		println(err)
	} else {
		buffer, IsAPausedGame, IsPlayingSolo = backEnd.GameSelection(Player)
	}

	return buffer, ServerPublicKey
}
