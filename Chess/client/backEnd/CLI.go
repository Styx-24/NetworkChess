package backEnd

import (
	"TP/structs"

	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func PlayerMenu() structs.User {
	var valid = false
	var option = 0
	var player structs.User
	var err error

	for !valid {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print(" 1. Choose a player \n 2. Create a new player \n >> ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ReplaceAll(text, "\n", "")

		option, err = strconv.Atoi(text)
		if err != nil || option < 1 || option > 2 {
			fmt.Println("Please enter a valid number")
		} else {
			valid = true
		}
	}

	switch option {
	case 1:
		player = LoadPlayer()

	case 2:
		player = CreatePlayer()
	}

	return player

}

func CreatePlayer() structs.User {

	var player structs.User

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name >> ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", "")
	player.Name = text

	fmt.Print("Enter your last name >> ")
	text, _ = reader.ReadString('\n')
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "\n", "")
	player.LastName = text

	player.Id = uuid.New()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	player.PrivateKey = *key
	if err != nil {
		log.Fatalf("Failed to generate RSA key: %v", err)
	}

	player.PublicKey = key.PublicKey

	InsertNewUser(player)

	return player
}

func LoadPlayer() structs.User {
	var list []structs.User
	var option = 0
	var valid = false
	var player structs.User

	list = GetUsers()

	for !valid {
		reader := bufio.NewReader(os.Stdin)
		print("Choose a player \n 0: quit \n")
		for i := 0; i < len(list); i++ {
			println(strconv.Itoa(i+1) + ": " + list[i].Name + " " + list[i].LastName)
		}
		print(">>")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ReplaceAll(text, "\n", "")

		response, err := strconv.Atoi(text)

		if err != nil || response < 0 || response > len(list) {
			fmt.Println("Please enter a valid number")
		} else {
			valid = true

			option = response - 1
		}
	}

	if option != -1 {
		key, err := rsa.GenerateKey(rand.Reader, 2048)
		list[option].PrivateKey = *key
		if err != nil {
			print(err)
		}

		list[option].PublicKey = key.PublicKey
		player = list[option]
	} else {
		player = PlayerMenu()
	}

	return player
}

func GameSelection(player structs.User) ([]byte, bool, bool) {

	option := 0
	valid := false
	buffer := make([]byte, 0)
	var err error
	var IsAPausedGame = false
	var IsPlayingSolo = false

	for !valid {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(" 1. Start game as host \n 2. Look for a game \n 3. Game against AI \n 4. Get paused games \n >> ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ReplaceAll(text, "\n", "")

		response, err := strconv.Atoi(text)

		if err != nil || response < 1 || response > 4 {
			fmt.Println("Please enter a valid number")
		} else {
			valid = true

			option = response
		}
	}

	switch option {
	case 1:
		{
			var GameRequest structs.GameRequest
			GameRequest.PlayerId = player.Id
			GameRequest.OponentId = player.Id

			buffer, err = GameRequest.Encode(player.PrivateKey)
			if err != nil {
				println(err)
			}

		}
	case 2:
		{

			var matchMakingRequest structs.MatchMakingRequest

			matchMakingRequest.ID = player.Id
			matchMakingRequest.IsAPausedGame = false
			buffer, err = matchMakingRequest.Encode(player.PrivateKey)
			if err != nil {
				println(err)
			}
		}
	case 3:
		{
			var GameRequest structs.GameRequest
			GameRequest.PlayerId = player.Id
			GameRequest.OponentId = uuid.Nil
			GameRequest.GameId = uuid.Nil
			IsPlayingSolo = true
			buffer, err = GameRequest.Encode(player.PrivateKey)
			if err != nil {
				println(err)
			}
		}
	case 4:
		{
			var matchMakingRequest structs.MatchMakingRequest

			IsAPausedGame = true
			matchMakingRequest.ID = player.Id
			matchMakingRequest.IsAPausedGame = true
			buffer, err = matchMakingRequest.Encode(player.PrivateKey)
			if err != nil {
				println(err)
			}
		}
	}

	return buffer, IsAPausedGame, IsPlayingSolo
}

func OpponentSelection(matchMakingResponse structs.MatchMakingResponse) int {
	var option = 0
	var valid = false

	for !valid {
		reader := bufio.NewReader(os.Stdin)

		print("Choose a player \n 0: quit \n")
		for i := 0; i < len(matchMakingResponse.Names); i++ {
			println(strconv.Itoa(i+1) + ": " + matchMakingResponse.Names[i])
		}
		print(">>")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ReplaceAll(text, "\n", "")

		response, err := strconv.Atoi(text)

		if err != nil || response < 0 || response > len(matchMakingResponse.Names) {
			fmt.Println("Please enter a valid number")
		} else {
			valid = true
			option = response
		}
	}

	return option
}

func SelectMove(player structs.User, gameId uuid.UUID, encryptionKey []byte, isPlayingSolo bool) []byte {
	var buffer []byte
	var err error

	reader := bufio.NewReader(os.Stdin)
	option := 0
	text := ""
	isValid := false

	for !isValid {
		fmt.Print("1: List of valid moves \n2: Best move \n3: Ask for a null game \n4: Pause game \n\nMake your move >> ")
		text, _ = reader.ReadString('\n')
		text = strings.ReplaceAll(text, "\n", "")
		text = strings.ReplaceAll(text, "\r", "")

		option, err = strconv.Atoi(text)
		if err == nil {
			if option > 0 && option < 5 {
				isValid = true
			} else {
				println("Please enter a valid number")
			}

		} else {
			isValid = true
		}
	}

	switch option {
	case 0:
		{
			var request structs.ActionRequest

			request.PlayerId = player.Id
			request.GameId = gameId
			request.Move = text

			buffer, err = request.Encode(player.PrivateKey, encryptionKey)
			if err != nil {
				println(err.Error())
			}
		}
	case 3:
		{
			if !isPlayingSolo {
				var request structs.DrawRequest
				request.GameId = gameId
				request.PlayerId = player.Id
				request.Message = "The opponent wants to null the game, do you accept?"
				buffer, err = request.Encode(player.PrivateKey)
				if err != nil {
					println(err.Error())
				}
			} else {
				println("Option not available in a game against the AI")
				return SelectMove(player, gameId, encryptionKey, isPlayingSolo)
			}

		}
	case 4:
		{

			var request structs.PauseRequest
			request.GameId = gameId
			request.PlayerId = player.Id
			request.Message = "The opponent wants to pause the game, do you accept?"
			buffer, err = request.Encode(player.PrivateKey)
			if err != nil {
				println(err.Error())
			}

		}
	default:
		{
			var request structs.InfoRequest

			request.PlayerId = player.Id
			request.GameId = gameId
			request.ValidMoves = option == 1

			buffer, err = request.Encode(player.PrivateKey, encryptionKey)
			if err != nil {
				println(err.Error())
			}
		}
	}

	return buffer
}

func ComfirmationPromt() int {

	option := 0
	valid := false

	for !valid {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(" 1. Accept \n 2. Refuse \n >> ")

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		text = strings.ReplaceAll(text, "\n", "")

		response, err := strconv.Atoi(text)

		if err != nil || response < 1 || response > 2 {
			fmt.Println("Please enter a valid number")
		} else {
			valid = true

			option = response
		}
	}

	return option

}
