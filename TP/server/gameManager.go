package server

import (
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

func GenerateGame() *chess.Game {
	return chess.NewGame()
}

func LoadGame(fenStr string) *chess.Game {
	fen, _ := chess.FEN(fenStr)
	return chess.NewGame(fen)

}

func Move(game *chess.Game, movestr string) error {
	if err := game.MoveStr(movestr); err != nil {
		return err
	}
	return nil
}

func AIMove(game *chess.Game) {
	eng, err := uci.New("C:\\Users\\atori\\Downloads\\stockfish\\stockfish.exe")
	if err != nil {
		panic(err)
	}
	defer eng.Close()

	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		panic(err)
	}

	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: time.Second / 100}
	if err := eng.Run(cmdPos, cmdGo); err != nil {
		panic(err)
	}
	move := eng.SearchResults().BestMove
	if err := game.Move(move); err != nil {
		panic(err)
	}

}

func CheckVictory(game *chess.Game) (int, string) {

	outcome := 0
	method := "none"

	switch game.Outcome() {
	case chess.BlackWon:
		outcome = 1
	case chess.WhiteWon:
		outcome = 2
	case chess.Draw:
		outcome = 3
	}

	if outcome != 0 {
		method = game.Method().String()
	}

	return outcome, method
}

func GetValidMoves(game *chess.Game) string {
	moves := ""
	moveList := game.ValidMoves()

	for i := range moveList {
		moves = moves + moveList[i].String() + "\n"
	}

	return moves
}

func GetBestMove(game *chess.Game) (string, error) {

	eng, err := uci.New("C:\\Users\\atori\\Downloads\\stockfish\\stockfish.exe")
	if err != nil {
		return "", err
	}
	defer eng.Close()

	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdPosition{game.Position(), game.Moves()}); err != nil {
		return "", err
	}

	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: time.Second / 100}

	if err := eng.Run(cmdPos, cmdGo); err != nil {
		return "", err
	}

	return eng.SearchResults().BestMove.String(), nil

}
