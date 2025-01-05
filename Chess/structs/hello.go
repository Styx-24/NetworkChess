package structs

import (
	"TP/protocols"
	"TP/utils"
	"crypto/rsa"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"

	"github.com/google/uuid"
)

type HelloRequest struct {
	Name     string
	LastName string
	Id       uuid.UUID
	Key      rsa.PublicKey
}

func (h HelloRequest) Decode(data []byte, player *User) error {

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
				player.Id, err = uuid.FromBytes(value)
				if err != nil {
					return err
				}
			}
		case protocols.String:
			{
				if player.Name == "" {
					player.Name = string(value)
				} else {
					player.LastName = string(value)
				}
			}
		case protocols.Byte:
			{
				block, _ := pem.Decode(value)
				if block == nil || block.Type != "PUBLIC KEY" {
					return err
				}

				pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					return err
				}

				rsaPubKey, ok := pubKey.(*rsa.PublicKey)
				if !ok {
					return err
				}

				player.PublicKey = *rsaPubKey

			}
		}

	}

	return nil

}

func (h HelloRequest) Encode(player User) ([]byte, error) {
	var buffer = make([]byte, 0)

	utils.AddStringToBuffer(player.Name, &buffer)

	utils.AddStringToBuffer(player.LastName, &buffer)

	utils.AddIdToBuffer(protocols.UUIDClient, player.Id, &buffer)

	utils.AddPublicKeyToBuffer(&player.PublicKey, &buffer)

	return utils.MakeFinalBuffer(protocols.HelloRequest, buffer), nil

}

type HelloResponse struct {
	Key rsa.PublicKey
}

func (h HelloResponse) Decode(data []byte) (rsa.PublicKey, error) {

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
		case 13:
			{
				block, _ := pem.Decode(value)
				if block == nil || block.Type != "PUBLIC KEY" {
					return rsa.PublicKey{}, err
				}

				pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
				if err != nil {
					return rsa.PublicKey{}, err
				}

				rsaPubKey, ok := pubKey.(*rsa.PublicKey)
				if !ok {
					return rsa.PublicKey{}, err
				}

				return *rsaPubKey, nil

			}
		}

	}

	return rsa.PublicKey{}, nil

}

func (h HelloResponse) Encode(key rsa.PublicKey) ([]byte, error) {
	var buffer = make([]byte, 0)

	utils.AddPublicKeyToBuffer(&key, &buffer)

	return utils.MakeFinalBuffer(protocols.HelloResponse, buffer), nil

}
