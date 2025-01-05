package server

import (
	"TP/protocols"
	"TP/server/TLV"
	"TP/server/backEnd"
	"TP/structs"
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/google/uuid"
)

var stockfishPath = ""

func Server() {

	var port = ""
	TLV.Games = make(map[uuid.UUID]structs.Game, 0)
	TLV.Players = make(map[uuid.UUID]structs.User, 0)
	TLV.GameMatchMaking = make(map[uuid.UUID]structs.Game, 0)

	genKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		println(err)
	}
	TLV.PrivateKey = *genKey

	err = backEnd.DbCreation()
	if err != nil {
		fmt.Println(err)
		return
	}

	port, stockfishPath = GetConfig()

	backEnd.StockfishPath = stockfishPath

	PORT := port
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {

		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go HandleConnection(c)

	}

}

func HandleConnection(c net.Conn) {

	println("new connexion")

	for {
		buffer := make([]byte, 2048)
		response := make([]byte, 0)
		pos := 0

		n, err := bufio.NewReader(c).Read(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}

		for pos < n {
			t := buffer[pos]
			pos++

			l := binary.LittleEndian.Uint16(buffer[pos:(pos + 2)])
			pos += 2
			value := buffer[pos : pos+int(l)]
			pos += int(l)

			switch t {
			case protocols.HelloRequest:
				{
					response = TLV.HelloRequest(value)
				}
			case protocols.GameRequest:
				{
					response = TLV.GameRequest(value, c)
				}
			case protocols.GameComfirmationResponse:
				{
					response = TLV.GameComfirmationResponse(value)
				}
			case protocols.DrawRequest:
				{
					response = TLV.DrawRequest(value)
				}
			case protocols.DrawResponse:
				{
					response = TLV.DrawResponse(value)
				}
			case protocols.PauseRequest:
				{
					response = TLV.PauseRequest(value)
				}
			case protocols.PauseResponse:
				{
					response = TLV.PauseResponse(value)
				}
			case protocols.ActionRequest:
				{
					response = TLV.ActionRequest(value)
				}
			case protocols.MatchMakingRequest:
				{
					response = TLV.MatchMakingRequest(value)
				}
			case protocols.InfoRequest:
				{
					response = TLV.InfoRequest(value)

				}
			}

		}

		if len(response) != 0 {
			c.Write(response)
		}
	}
}
