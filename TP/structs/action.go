package structs

import (
	"TP/protocols"
	"TP/utils"
	"crypto/rsa"
	"encoding/binary"

	"github.com/google/uuid"
)

type ActionRequest struct {
	PlayerId  uuid.UUID
	GameId    uuid.UUID
	Move      string
	Signature []byte
}

func (h ActionRequest) Decode(data []byte, gameList map[uuid.UUID]Game) (ActionRequest, []byte, error) {
	var tmp ActionRequest

	var buffer = make([]byte, 0)
	pos := 0
	var err error
	var key = make([]byte, 0)

	for (pos < len(data)) && tmp.GameId == uuid.Nil {

		t := data[pos]
		pos++

		l := binary.LittleEndian.Uint16(data[pos:(pos + 2)])
		pos += 2
		value := data[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.UUIDClient:
			{

				tmp.PlayerId, err = uuid.FromBytes(value)
				if err != nil {
					return tmp, make([]byte, 0), err
				}

			}
		case protocols.UUIDGame:
			{

				tmp.GameId, err = uuid.FromBytes(value)
				if err != nil {
					return tmp, make([]byte, 0), err
				}

			}
		}

	}

	key = gameList[tmp.GameId].EncryptionKey
	decryptedData := utils.Decrypt(key, data[pos:])
	pos = 0

	for pos < len(decryptedData) {

		t := decryptedData[pos]
		pos++

		l := binary.LittleEndian.Uint16(decryptedData[pos:(pos + 2)])
		pos += 2
		value := decryptedData[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.String:
			{

				tmp.Move = string(value)
				if err != nil {
					return tmp, make([]byte, 0), err
				} else {
					utils.AddStringToBuffer(tmp.Move, &buffer)
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

func (h ActionRequest) Encode(privateKey rsa.PrivateKey, key []byte) ([]byte, error) {

	var buffer = make([]byte, 0)

	utils.AddStringToBuffer(h.Move, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	encryptedBuffer := utils.Encrypt(key, buffer)

	buffer = make([]byte, 0)

	utils.AddIdToBuffer(protocols.UUIDClient, h.PlayerId, &buffer)

	utils.AddIdToBuffer(protocols.UUIDGame, h.GameId, &buffer)

	var finalBuffer = make([]byte, 0)
	finalBuffer = append(finalBuffer, byte(protocols.ActionRequest))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(encryptedBuffer)+len(buffer)))
	finalBuffer = append(finalBuffer, l...)
	finalBuffer = append(finalBuffer, buffer...)
	finalBuffer = append(finalBuffer, encryptedBuffer...)

	return finalBuffer, nil

}

type ActionResponse struct {
	MoveWasValid bool
	GameHasEnded bool
	Message      string
	TurnOf       int
	Signature    []byte
}

func (h ActionResponse) Decode(data []byte, key []byte) (ActionResponse, []byte, error) {
	var tmp ActionResponse

	var buffer = make([]byte, 0)
	pos := 0

	moveWasValideAssinged := false
	gameHasEndedAssinged := false

	decryptedData := utils.Decrypt(key, data)

	for pos < len(decryptedData) {

		t := decryptedData[pos]
		pos++

		l := binary.LittleEndian.Uint16(decryptedData[pos:(pos + 2)])
		pos += 2
		value := decryptedData[pos : pos+int(l)]
		pos += int(l)

		switch t {
		case protocols.String:
			{

				tmp.Message = string(value)

				utils.AddStringToBuffer(string(value), &buffer)

			}
		case protocols.Int:
			{

				if !moveWasValideAssinged {
					tmp.MoveWasValid = value[0] == 1
					moveWasValideAssinged = true
				} else if !gameHasEndedAssinged {
					tmp.GameHasEnded = value[0] == 1
					gameHasEndedAssinged = true
				} else {
					tmp.TurnOf = int(value[0])
				}

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

func (h ActionResponse) Encode(privateKey rsa.PrivateKey, key []byte) ([]byte, error) {

	var buffer = make([]byte, 0)

	intMoveWasValid := 0
	intGameHasEnded := 0

	if h.MoveWasValid {
		intMoveWasValid = 1
	}

	if h.GameHasEnded {
		intGameHasEnded = 1
	}

	utils.AddIntToBuffer(intMoveWasValid, &buffer)

	utils.AddIntToBuffer(intGameHasEnded, &buffer)

	utils.AddStringToBuffer(h.Message, &buffer)

	utils.AddIntToBuffer(h.TurnOf, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	encryptedBuffer := utils.Encrypt(key, buffer)

	return utils.MakeFinalBuffer(protocols.ActionResponse, []byte(encryptedBuffer)), nil

}
