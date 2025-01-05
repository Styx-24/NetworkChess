package structs

import (
	"TP/protocols"
	"TP/utils"
	"crypto/rsa"
	"encoding/binary"

	"github.com/google/uuid"
)

type GameRequest struct {
	PlayerId  uuid.UUID
	OponentId uuid.UUID
	GameId    uuid.UUID
	Signature []byte
}

func (h GameRequest) Decode(data []byte) (GameRequest, []byte, error) {
	var tmp GameRequest

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
				if tmp.PlayerId == uuid.Nil {

					tmp.PlayerId, err = uuid.FromBytes(value)
					if err != nil {
						return h, make([]byte, 0), err
					} else {
						utils.AddIdToBuffer(protocols.UUIDClient, tmp.PlayerId, &buffer)
					}

				} else if tmp.OponentId == uuid.Nil {

					tmp.OponentId, err = uuid.FromBytes(value)
					if err != nil {
						return h, make([]byte, 0), err
					} else {
						utils.AddIdToBuffer(protocols.UUIDClient, tmp.OponentId, &buffer)
					}

				}
			}
		case protocols.UUIDGame:
			{
				tmp.GameId, err = uuid.FromBytes(value)
				if err != nil {
					return h, make([]byte, 0), err
				} else {
					utils.AddIdToBuffer(protocols.UUIDGame, tmp.GameId, &buffer)
				}
			}
		case protocols.Signature:
			{
				tmp.Signature = value
			}
		}

	}

	return tmp, buffer, nil

}

func (h GameRequest) Encode(privateKey rsa.PrivateKey) ([]byte, error) {
	var buffer = make([]byte, 0)

	utils.AddIdToBuffer(protocols.UUIDClient, h.PlayerId, &buffer)

	utils.AddIdToBuffer(protocols.UUIDClient, h.OponentId, &buffer)

	utils.AddIdToBuffer(protocols.UUIDGame, h.GameId, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	return utils.MakeFinalBuffer(protocols.GameRequest, buffer), nil

}

type GameResponse struct {
	GameId        uuid.UUID
	Status        string
	Team          int
	TurnOf        int
	EncryptionKey []byte
	Signature     []byte
}

func (h GameResponse) Decode(data []byte) (GameResponse, []byte, error) {
	var tmp GameResponse
	var buffer = make([]byte, 0)
	pos := 0
	var err error
	var teamIsSet = false

	for pos < len(data) {

		t := data[pos]
		pos++

		l := binary.LittleEndian.Uint16(data[pos:(pos + 2)])
		pos += 2
		value := data[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.String:
			{
				tmp.Status = string(value)
				utils.AddStringToBuffer(tmp.Status, &buffer)
			}
		case protocols.Int:
			{
				if !teamIsSet {
					tmp.Team = int(value[0])
					teamIsSet = true
				} else {
					tmp.TurnOf = int(value[0])
				}

				utils.AddIntToBuffer(int(value[0]), &buffer)
			}
		case protocols.Byte:
			{
				tmp.EncryptionKey = value
				utils.AddBytesToBuffer(value, &buffer)
			}
		case protocols.UUIDGame:
			{

				tmp.GameId, err = uuid.FromBytes(value)
				if err != nil {
					return h, make([]byte, 0), err
				} else {
					utils.AddIdToBuffer(2, tmp.GameId, &buffer)
				}

			}
		case protocols.Signature:
			{
				tmp.Signature = value
			}
		}

	}

	return tmp, buffer, nil

}

func (h GameResponse) Encode(privateKey rsa.PrivateKey) ([]byte, error) {
	var buffer = make([]byte, 0)

	utils.AddIdToBuffer(protocols.UUIDGame, h.GameId, &buffer)

	utils.AddStringToBuffer(h.Status, &buffer)

	utils.AddIntToBuffer(h.Team, &buffer)

	utils.AddIntToBuffer(h.TurnOf, &buffer)

	utils.AddBytesToBuffer(h.EncryptionKey, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	return utils.MakeFinalBuffer(protocols.GameResponse, buffer), nil

}
