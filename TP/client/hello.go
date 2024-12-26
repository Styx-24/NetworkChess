package client

import (
	"TP/structs"
	"crypto/rsa"
)

func HelloResponse(value []byte) ([]byte, rsa.PublicKey) {
	var buffer []byte

	var helloResponse structs.HelloResponse
	serverPublicKey, err := helloResponse.Decode(value)
	serverPublicKey.Size()
	if err != nil {
		println(err)
	} else {
		println("Connexion établis avec success")
		buffer, IsAPausedGame = GameSelection(player)
	}

	return buffer, serverPublicKey
}
