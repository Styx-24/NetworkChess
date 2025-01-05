package client

import (
	"TP/client/TLV"
	"TP/client/backEnd"
	"TP/protocols"
	"TP/structs"
	"bufio"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/google/uuid"
)

func Client() {

	CONNECT := GetServerAddress()
	c, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	TLV.Player.Id = uuid.New()
	TLV.Player.Name = "Michel"
	TLV.Player.LastName = "Jacob"

	err = backEnd.DbCreation()
	if err != nil {
		fmt.Println(err)
		return
	}

	TLV.Player = backEnd.PlayerMenu()

	buffer, err := structs.HelloRequest.Encode(structs.HelloRequest{}, TLV.Player)
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
					buffer, TLV.ServerPublicKey = TLV.HelloResponse(value)
				}
			case protocols.GameResponse:
				{
					buffer = TLV.GameResponse(value)
				}

			case protocols.ActionResponse:
				{
					buffer = TLV.ActionResponse(value)
				}
			case protocols.MatchMakingResponse:
				{
					buffer = TLV.MatchMakingResponse(value)
				}
			case protocols.InfoResponse:
				{
					buffer = TLV.InfoResponse(value)
				}
			case protocols.GameComfirmationRequest:
				{
					buffer = TLV.GameComfirmationRequest(value)
				}
			case protocols.GameComfirmationResponse:
				{
					buffer = TLV.GameComfirmationResponse(value)
				}
			case protocols.DrawRequest:
				{
					buffer = TLV.DrawRequest(value)
				}
			case protocols.PauseRequest:
				{
					buffer = TLV.PauseRequest(value)
				}
			}

			if len(buffer) > 0 {
				c.Write(buffer)
			}

		}

	}
}
