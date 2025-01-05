package structs

import (
	"TP/protocols"
	"TP/utils"
	"crypto/rsa"
	"encoding/binary"

	"github.com/google/uuid"
)

type PauseRequest struct {
	PlayerId  uuid.UUID
	GameId    uuid.UUID
	Message   string
	Signature []byte
}

func (h PauseRequest) Decode(data []byte) (PauseRequest, []byte, error) {
	var tmp PauseRequest

	var buffer = make([]byte, 0)
	var err error
	pos := 0

	for pos < len(data) {

		t := data[pos]
		pos++

		l := binary.LittleEndian.Uint16(data[pos:(pos + 2)])
		pos += 2
		value := data[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.UUIDGame:
			{
				tmp.GameId, err = uuid.FromBytes(value)
				if err != nil {
					return h, make([]byte, 0), err
				} else {
					utils.AddIdToBuffer(protocols.UUIDGame, tmp.GameId, &buffer)
				}
			}
		case protocols.UUIDClient:
			{
				tmp.PlayerId, err = uuid.FromBytes(value)
				if err != nil {
					return h, make([]byte, 0), err
				} else {
					utils.AddIdToBuffer(protocols.UUIDClient, tmp.PlayerId, &buffer)
				}
			}
		case protocols.String:
			{
				tmp.Message = string(value)
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

func (h PauseRequest) Encode(privateKey rsa.PrivateKey) ([]byte, error) {
	var buffer = make([]byte, 0)

	utils.AddIdToBuffer(protocols.UUIDGame, h.GameId, &buffer)
	utils.AddIdToBuffer(protocols.UUIDClient, h.PlayerId, &buffer)

	utils.AddStringToBuffer(h.Message, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	return utils.MakeFinalBuffer(protocols.PauseRequest, buffer), nil

}

type PauseResponse struct {
	PlayerId  uuid.UUID
	GameId    uuid.UUID
	Answer    bool
	Signature []byte
}

func (h PauseResponse) Decode(data []byte) (PauseResponse, []byte, error) {
	var tmp PauseResponse
	var buffer = make([]byte, 0)
	var err error
	pos := 0

	for pos < len(data) {

		t := data[pos]
		pos++

		l := binary.LittleEndian.Uint16(data[pos:(pos + 2)])
		pos += 2
		value := data[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.UUIDGame:
			{
				tmp.GameId, err = uuid.FromBytes(value)
				if err != nil {
					return h, make([]byte, 0), err
				} else {
					utils.AddIdToBuffer(protocols.UUIDGame, tmp.GameId, &buffer)
				}
			}
		case protocols.UUIDClient:
			{
				tmp.PlayerId, err = uuid.FromBytes(value)
				if err != nil {
					return h, make([]byte, 0), err
				} else {
					utils.AddIdToBuffer(protocols.UUIDClient, tmp.PlayerId, &buffer)
				}
			}
		case protocols.Int:
			{
				tmp.Answer = value[0] == 1
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

func (h PauseResponse) Encode(privateKey rsa.PrivateKey) ([]byte, error) {
	var buffer = make([]byte, 0)

	intAnswer := 0

	if h.Answer {
		intAnswer = 1
	}

	utils.AddIdToBuffer(protocols.UUIDGame, h.GameId, &buffer)
	utils.AddIdToBuffer(protocols.UUIDClient, h.PlayerId, &buffer)
	utils.AddIntToBuffer(intAnswer, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	return utils.MakeFinalBuffer(protocols.PauseResponse, buffer), nil

}
