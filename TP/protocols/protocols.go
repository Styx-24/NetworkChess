package protocols

// Request for the first identification to the server where the server and client exchanges public keys
const HelloRequest = 0
const HelloResponse = 100

// data types for request encoding/decoding
const UUIDClient = 1
const UUIDGame = 2
const Signature = 3
const String = 11
const Int = 12
const Byte = 13

// Request to start a game vs another player or vs the server
const GameRequest = 30
const GameResponse = 130

// Request to make a move in an active chess game
const ActionRequest = 40
const ActionResponse = 140

// Request to get all the waiting players on the server
const MatchMakingRequest = 110
const MatchMakingResponse = 120

// Request to get the valid moves and best moves
const InfoRequest = 115
const InfoResponse = 125

// Request to get if the user has accepted the match or not
const GameComfirmationRequest = 35
const GameComfirmationResponse = 135

// Request to get if the user has accepted the draw request or not
const DrawRequest = 36
const DrawResponse = 136

// Request to get if the user has accepted the pause request or not
const PauseRequest = 37
const PauseResponse = 137
