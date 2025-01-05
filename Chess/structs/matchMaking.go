package structs

import (
	"TP/protocols"
	"TP/utils"
	"crypto/rsa"
	"encoding/binary"

	"github.com/google/uuid"
)

type MatchMakingRequest struct {
	ID            uuid.UUID
	IsAPausedGame bool
	Signature     []byte
}

func (h MatchMakingRequest) Decode(data []byte) (MatchMakingRequest, []byte, error) {
	var tmp MatchMakingRequest
	var buffer = make([]byte, 0)
	pos := 0
	var err error

	for pos < len(data) {

		t := data[pos]
		pos++

		l := binary.LittleEndian.Uint16(data[pos:(pos + 2)])
		pos += 2
		value := data[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.UUIDClient:
			{
				tmp.ID, err = uuid.FromBytes(value)
				if err != nil {
					print(err)
				}

				utils.AddIdToBuffer(1, tmp.ID, &buffer)
			}
		case protocols.Int:
			{
				tmp.IsAPausedGame = value[0] == 1
				utils.AddIntToBuffer(int(value[0]), &buffer)
			}
		case protocols.Signature:
			{
				tmp.Signature = value
			}
		}

	}

	return tmp, buffer, nil

}

func (h MatchMakingRequest) Encode(privateKey rsa.PrivateKey) ([]byte, error) {
	var buffer = make([]byte, 0)
	var intIsAPausedGame = 0

	if h.IsAPausedGame {
		intIsAPausedGame = 1
	}

	utils.AddIdToBuffer(protocols.UUIDClient, h.ID, &buffer)
	utils.AddIntToBuffer(intIsAPausedGame, &buffer)
	utils.SignBuffer(privateKey, &buffer)

	return utils.MakeFinalBuffer(protocols.MatchMakingRequest, buffer), nil

}

type MatchMakingResponse struct {
	IDs       []uuid.UUID
	Names     []string
	Signature []byte
}

func (h MatchMakingResponse) Decode(data []byte) (MatchMakingResponse, []byte, error) {
	var tmp MatchMakingResponse
	tmp.IDs = make([]uuid.UUID, 0)
	tmp.Names = make([]string, 0)

	var buffer = make([]byte, 0)
	pos := 0

	for pos < len(data) {

		t := data[pos]
		pos++

		l := binary.LittleEndian.Uint16(data[pos:(pos + 2)])
		pos += 2
		value := data[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.UUIDClient:
			{
				id, err := uuid.FromBytes(value)
				if err != nil {
					print(err)
				}
				tmp.IDs = append(tmp.IDs, id)

				utils.AddIdToBuffer(1, id, &buffer)
			}
		case protocols.String:
			{
				tmp.Names = append(tmp.Names, string(value))

				utils.AddStringToBuffer(string(value), &buffer)
			}
		case protocols.Signature:
			{
				tmp.Signature = value
			}
		}

	}

	return tmp, buffer, nil

}

func (h MatchMakingResponse) Encode(privateKey rsa.PrivateKey) ([]byte, error) {
	var buffer = make([]byte, 0)

	for i := 0; i < len(h.IDs); i++ {
		utils.AddIdToBuffer(protocols.UUIDClient, h.IDs[i], &buffer)
		utils.AddStringToBuffer(h.Names[i], &buffer)
	}

	utils.SignBuffer(privateKey, &buffer)

	return utils.MakeFinalBuffer(protocols.MatchMakingResponse, buffer), nil

}
