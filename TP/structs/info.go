package structs

import (
	"TP/protocols"
	"TP/utils"
	"crypto/rsa"
	"encoding/binary"

	"github.com/google/uuid"
)

type InfoRequest struct {
	GameId     uuid.UUID
	PlayerId   uuid.UUID
	ValidMoves bool
	Signature  []byte
}

func (h InfoRequest) Decode(data []byte, gameList map[uuid.UUID]Game) (InfoRequest, []byte, error) {
	var tmp InfoRequest

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
		case protocols.Int:
			{

				tmp.ValidMoves = value[0] == 1

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

func (h InfoRequest) Encode(privateKey rsa.PrivateKey, key []byte) ([]byte, error) {

	var buffer = make([]byte, 0)

	validMoves := 0

	if h.ValidMoves {
		validMoves = 1
	}

	utils.AddIntToBuffer(validMoves, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	encryptedBuffer := utils.Encrypt(key, buffer)

	buffer = make([]byte, 0)

	utils.AddIdToBuffer(protocols.UUIDClient, h.PlayerId, &buffer)
	utils.AddIdToBuffer(protocols.UUIDGame, h.GameId, &buffer)

	var finalBuffer = make([]byte, 0)
	finalBuffer = append(finalBuffer, byte(protocols.InfoRequest))
	l := make([]byte, 2)
	binary.LittleEndian.PutUint16(l, uint16(len(encryptedBuffer)+len(buffer)))
	finalBuffer = append(finalBuffer, l...)
	finalBuffer = append(finalBuffer, buffer...)
	finalBuffer = append(finalBuffer, encryptedBuffer...)

	return finalBuffer, nil

}

type InfoResponse struct {
	Move      string
	Signature []byte
}

func (h InfoResponse) Decode(data []byte, key []byte) (InfoResponse, []byte, error) {
	var tmp InfoResponse

	var buffer = make([]byte, 0)
	pos := 0

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
				tmp.Move = string(value)

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

func (h InfoResponse) Encode(privateKey rsa.PrivateKey, key []byte) ([]byte, error) {

	var buffer = make([]byte, 0)

	utils.AddStringToBuffer(h.Move, &buffer)

	utils.SignBuffer(privateKey, &buffer)

	encryptedBuffer := utils.Encrypt(key, buffer)

	return utils.MakeFinalBuffer(protocols.InfoResponse, []byte(encryptedBuffer)), nil

}
