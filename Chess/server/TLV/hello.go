package TLV

import (
	"TP/structs"
	"fmt"
)

func HelloRequest(value []byte) []byte {
	var HelloRequest structs.HelloRequest
	var HelloResponse structs.HelloResponse
	var player structs.User
	var response []byte

	err := HelloRequest.Decode(value, &player)
	if err != nil {
		println(err.Error())
	} else {
		Players[player.Id] = player
		fmt.Println("player : " + player.Name + " " + player.LastName + " added with ID " + player.Id.String())
		response, _ = HelloResponse.Encode(PrivateKey.PublicKey)
	}

	return response
}
