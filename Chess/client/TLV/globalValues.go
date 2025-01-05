package TLV

import (
	"TP/structs"
	"crypto/rsa"

	"github.com/google/uuid"
)

var Player structs.User
var ServerPublicKey rsa.PublicKey
var EncryptionKey []byte
var GameId uuid.UUID
var Team int
var IsAPausedGame = false
var IsPlayingSolo = false
