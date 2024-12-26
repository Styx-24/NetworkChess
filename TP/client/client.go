package client

import (
	"TP/protocols"
	"TP/structs"

	"bufio"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/google/uuid"
)

var player structs.User
var serverPublicKey rsa.PublicKey
var encryptionKey []byte
var gameId uuid.UUID
var team int
var IsAPausedGame = false
var isPlayingSolo = false

func Client() {

	CONNECT := GetServerAddress()
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	player.Id = uuid.New()
	player.Name = "Michel"
	player.LastName = "Jacob"

	err = DbCreation()
	if err != nil {
		fmt.Println(err)
		return
	}

	player = PlayerMenu()

	buffer, err := structs.HelloRequest.Encode(structs.HelloRequest{}, player)
	if err != nil {
		print(err)
	}

	c.Write(buffer)

	for {
		response := make([]byte, 2048)
		buffer = make([]byte, 0)

		n, err := bufio.NewReader(c).Read(response)
		if err != nil {
			fmt.Println(err)
		}

		pos := 0

		for pos < n {

			buffer = make([]byte, 0)
			t := response[pos]
			pos++

			l := binary.LittleEndian.Uint16(response[pos:(pos + 2)])
			pos += 2
			value := response[pos : pos+int(l)]
			pos += int(l)

			switch t {
			case protocols.HelloResponse:
				{
					buffer, serverPublicKey = HelloResponse(value)
				}
			case protocols.GameResponse:
				{
					buffer = GameResponse(value)
				}

			case protocols.ActionResponse:
				{
					buffer = ActionResponse(value)
				}
			case protocols.MatchMakingResponse:
				{
					buffer = MatchMakingResponse(value)
				}
			case protocols.InfoResponse:
				{
					buffer = InfoResponse(value)
				}
			case protocols.GameComfirmationRequest:
				{
					buffer = GameComfirmationRequest(value)
				}
			case protocols.GameComfirmationResponse:
				{
					buffer = GameComfirmationResponse(value)
				}
			case protocols.DrawRequest:
				{
					buffer = DrawRequest(value)
				}
			case protocols.PauseRequest:
				{
					buffer = PauseRequest(value)
				}
			}

			if len(buffer) > 0 {
				c.Write(buffer)
			}

		}

	}
}
