package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type InputCommand struct {
	client  *Client
	command []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type GameServer struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	command chan InputCommand

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
	// Game Data
	gamedata *GameData

	// Game Controller
	controller *GameController

	// UUID for next Client
	currentUUID int64
}

func newGameController() *GameController {
	return &GameController{}
}

func newGameServer() *GameServer {
	return &GameServer{
		broadcast:   make(chan []byte),
		command:     make(chan InputCommand),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		gamedata:    newGameData(),
		controller:  newGameController(),
		currentUUID: 100000,
	}
}
func (g *GameServer) getNextClientUUID() int64 {
	g.currentUUID += 1
	return g.currentUUID
}

func (g *GameServer) sendToAllClients(message []byte) {
	go func() { g.broadcast <- message }()
}

func (g *GameServer) clientCount() int {
	return len(g.clients)
}

func (g *GameServer) shutdown() {
	for client := range g.clients {

		client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

		//client.close <- true
		//select {
		//case client.close <- true:
		//default:
		//	close(client.send)
		//	delete(g.clients, client)
		//}
	}

	if g.clientCount() > 0 {
		time.Sleep(time.Millisecond * 2000)
	}
}
func (g *GameServer) run() {
	ticker := time.NewTicker(20 * time.Millisecond)
	saveticker := time.NewTicker(30 * time.Second)
	last := time.Now()
	defer func() {
		ticker.Stop()
		saveticker.Stop()
		g.controller.onEnd()
	}()

	g.controller.onInit(g)

	for {
		select {

		// Client Connect
		case client := <-g.register:
			g.clients[client] = true
			g.controller.onClientConnect(client)

		// Client Disconnect
		case client := <-g.unregister:

			if _, ok := g.clients[client]; ok {
				delete(g.clients, client)
				close(client.send)
			}

			g.controller.onClientDisconnect(client)

		// Timer
		case <-ticker.C:
			g.controller.onDataUpdate(time.Now().Sub(last))

		case <-saveticker.C:
			go g.gamedata.saveToFile()
		// Command Interpreter
		case input := <-g.command:
			g.controller.onParseMessage(input)

		// Broadcast message
		case message := <-g.broadcast:
			for client := range g.clients {
				select {
				case client.send <- string(message):
				default:
					close(client.send)
					delete(g.clients, client)
				}
			}
		}
	}
}
