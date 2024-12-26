package server

import (
	"TP/protocols"
	"TP/structs"
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/google/uuid"
)

var players = make(map[uuid.UUID]structs.User)
var games = make(map[uuid.UUID]structs.Game)
var gameMatchMaking = make(map[uuid.UUID]structs.Game)
var privateKey rsa.PrivateKey

func Server() {

	genKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		println(err)
	}
	privateKey = *genKey

	err = DbCreation()
	if err != nil {
		fmt.Println(err)
		return
	}

	PORT := ":8080"
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

	println("nouvelle connexion")

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
					response = HelloRequest(value)
				}
			case protocols.GameRequest:
				{
					response = GameRequest(value, c)
				}
			case protocols.GameComfirmationResponse:
				{
					response = GameComfirmationResponse(value)
				}
			case protocols.DrawRequest:
				{
					response = DrawRequest(value)
				}
			case protocols.DrawResponse:
				{
					response = DrawResponse(value)
				}
			case protocols.PauseRequest:
				{
					response = PauseRequest(value)
				}
			case protocols.PauseResponse:
				{
					response = PauseResponse(value)
				}
			case protocols.ActionRequest:
				{
					response = ActionRequest(value)
				}
			case protocols.MatchMakingRequest:
				{
					response = MatchMakingRequest(value)
				}
			case protocols.InfoRequest:
				{
					response = InfoRequest(value)

				}
			}

		}

		if len(response) != 0 {
			c.Write(response)
		}
	}
}
