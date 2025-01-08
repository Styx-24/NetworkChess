# Network Chess Game

This is a online chess game written in Go using the TCP protocol. It uses the [notnil go chess library](https://github.com/notnil/chess) and stockfish as backend.

## How to install

Download [the latest release](https://github.com/Styx-24/NetworkChess/releases/tag/1.0) and extract it.

## How to use

Before starting the program make sure to set the correct path for stockfish in the server config file located at Chess/server/config.txt this is not necessary if you only intend to use the client. you can also set the port that the server will listen to in that file(default is 8080). You can change the port that the client will try to connect to(default is localhost 8080) in the client config file at Chess/client/config.txt.

To start the program open a terminal at the root of the project and then type the command

```
go run main.go s
```


or 

```
go run main.go c
```
to start the program in server or client mode.

When in a game you can use normal coordinates(ex:a2a4) or algebraic chess notations(ex:a4) to input a move. When inputting a move in an active game, if the move captures another piece you will have to put a x between the two coordinates like this: e2xe4. Otherwise, the move will be seen as invalid by stockfish even if it is not written this way when sotckfish suggests the best move possible.

## Depandancies

[go](https://go.dev/doc/install)

[stockfish](https://stockfishchess.org/download/)
