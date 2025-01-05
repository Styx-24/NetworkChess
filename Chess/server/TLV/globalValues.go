package TLV

import (
	"TP/structs"
	"crypto/rsa"

	"github.com/google/uuid"
)

var Games map[uuid.UUID]structs.Game
var Players map[uuid.UUID]structs.User
var PrivateKey rsa.PrivateKey
var GameMatchMaking map[uuid.UUID]structs.Game
