package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"
)

func createClientDataUpdateCommand(client *Client) []byte {
	return CreateSimpleCommand(SIMPLE_COMMAND_UPDATE_OBJECT_DATA,
		client.data.clientUUID,
		"player",
		strconv.FormatInt(client.data.position.x, 10),
		strconv.FormatInt(client.data.position.y, 10),
		ParseColorToString(client.data.colorID),
	)
}
func createObjectAdd(client *Client) []byte {
	return CreateSimpleCommand(SIMPLE_COMMAND_UPDATE_OBJECT_DATA,
		client.data.clientUUID,
		"player",
		strconv.FormatInt(client.data.position.x, 10),
		strconv.FormatInt(client.data.position.y, 10),
		ParseColorToString(client.data.colorID),
	)
}

type GameController struct {
	server          *GameServer
	stateController GameStateController
}

func (g *GameController) changeState(state string) {

	if g.stateController != nil {
		g.stateController.onLeaveState()
	}

	switch state {
	case GAME_STATE_START:
		g.server.gamedata.blocks = []BlockData{}
		g.stateController = &GameStateStart{GameStateBase{controller: g}}
	case GAME_STATE_MAIN:
		g.stateController = &GameStateMain{GameStateBase{controller: g}}
	default:
		return
	}

	g.stateController.onEnterState()
}

func (g *GameController) onInit(server *GameServer) {
	log.Println("Controller: OnInit ")
	g.server = server
	g.changeState(GAME_STATE_MAIN)
}

func (g *GameController) onEnd() {
	log.Println("Controller: OnInit ")

	g.stateController.onLeaveState()
}

// Vorsicht bei Statuswechsel etc.... nur zum ausprobieren
func (g *GameController) writeTextToAll(text string, x int64, y int64, duration time.Duration) {

	message_text := CreateSimpleCommand(SIMPLE_COMMAND_LOG_MESSAGE, 0, text)
	g.server.sendToAllClients(message_text)

	/*g.server.gamedata.currentBlockUUID++
	cur_txt_id := g.server.gamedata.currentBlockUUID
	message_text := CreateSimpleCommand(SIMPLE_COMMAND_ADD_TEXT_DATA, cur_txt_id,
		strconv.FormatInt(x, 10),
		strconv.FormatInt(y, 10),
		text,
		"18px Arial",
		"black",
		"left")
	g.server.sendToAllClients(message_text)

	time.Sleep(duration)

	message_text_remove := CreateSimpleCommand(SIMPLE_COMMAND_REMOVE_TEXT_DATA, cur_txt_id)
	g.server.sendToAllClients(message_text_remove)*/
}

func (g *GameController) onClientConnect(client *Client) {
	log.Println("Controller: Client connected from ", client.conn.RemoteAddr())
	message := CreateSimpleCommand(SIMPLE_COMMAND_ADD_OBJECT_DATA,
		client.data.clientUUID,
		GOT_PLAYER,
		strconv.FormatInt(client.data.position.x, 10),
		strconv.FormatInt(client.data.position.y, 10),
		ParseColorToString(client.data.colorID),
	)

	g.server.sendToAllClients(message)

	go g.writeTextToAll("Client connected ID="+strconv.FormatInt(client.data.clientUUID, 10), 10, 380, time.Second*5)

}

func (g *GameController) onClientDisconnect(client *Client) {
	log.Println("Controller: Client disconnected from ", client.conn.RemoteAddr())
	message := CreateSimpleCommand(SIMPLE_COMMAND_REMOVE_OBJECT_DATA,
		client.data.clientUUID,
		GOT_PLAYER,
		strconv.FormatInt(client.data.position.x, 10),
		strconv.FormatInt(client.data.position.y, 10),
		ParseColorToString(client.data.colorID),
	)

	g.server.sendToAllClients(message)

	go g.writeTextToAll("Client disconnected ID="+strconv.FormatInt(client.data.clientUUID, 10), 10, 380, time.Second*5)

}

func (g *GameController) onDataUpdate(delta time.Duration) {
	g.stateController.onUpdate(delta)
}
func (g *GameController) sendConfigToClient(client *Client) {
	message := CreateSimpleCommand(SIMPLE_COMMAND_CONFIG, client.data.clientUUID,
		strconv.FormatInt(client.data.clientUUID, 10),
		strconv.FormatInt(CLIENT_STEP_SIZE, 10),
		strconv.FormatInt(CLIENT_STEP_SIZE, 10),
		strconv.FormatInt(client.data.token, 10),
	)

	client.send <- string(message)
}

func (g *GameController) sendAllDataToClient(client *Client) {

	//log.Println("Request all data: Object Count", len(g.server.gamedata.blocks))

	for y := int64(0); y < MAP_HEIGHT; y++ {
		for x := int64(0); x < MAP_WIDTH; x++ {
			//if len(g.server.gamedata.matrix[y][x].blockType) > 0 {
			if g.server.gamedata.matrix[y][x].colorID.A > 0 {
				message := CreateSimpleCommand(SIMPLE_COMMAND_ADD_OBJECT_DATA, y*MAP_WIDTH+x, g.server.gamedata.matrix[y][x].blockType,
					strconv.FormatInt(x*CLIENT_STEP_SIZE, 10),
					strconv.FormatInt(y*CLIENT_STEP_SIZE, 10),
					ParseColorToString(g.server.gamedata.matrix[y][x].colorID),
				)
				client.send <- string(message)
				//time.Sleep(1 * time.Millisecond)
			}
		}
	}

	/*for _, cur_object := range g.server.gamedata.blocks {
		message := CreateSimpleCommand(SIMPLE_COMMAND_ADD_OBJECT_DATA, cur_object.blockUUID, cur_object.blockType,
			strconv.FormatInt(cur_object.position.x, 10),
			strconv.FormatInt(cur_object.position.y, 10))
		client.send <- message
	}*/

	for cur_client := range g.server.clients {
		message := createClientDataUpdateCommand(cur_client)
		client.send <- string(message)
	}

}

func (g *GameController) onParseMessage(input InputCommand) {
	log.Println("Command from ", input.client.data.clientUUID, " : ", string(input.command))

	var jsonCommand JsonClientCommand

	err := json.Unmarshal(input.command, &jsonCommand)

	if err != nil {
		return
	}

	switch jsonCommand.CommandType {

	// Process Simple Command
	case COMMAND_TYPE_SIMPLE:
		var jsonSimpleCommand JsonClientSimpleCommand

		err := json.Unmarshal(jsonCommand.CommandData, &jsonSimpleCommand)

		if err != nil {
			log.Println("Invalid JSON data received")
			return
		}
		if jsonSimpleCommand.Id == input.client.data.token {
			if jsonSimpleCommand.Command == SIMPLE_COMMAND_KEY_PRESSED {
				g.stateController.onProcessInput(input.client, jsonSimpleCommand.Parameters)
			}

			if jsonSimpleCommand.Command == SIMPLE_COMMAND_SET_COLOR {
				g.stateController.onProcessCommand(input.client, jsonSimpleCommand.Command, jsonSimpleCommand.Parameters)
			}
		}
	// Process Uni Command
	case COMMAND_TYPE_UNI:
		var jsonUniCommand JsonClientUniCommand

		err := json.Unmarshal(jsonCommand.CommandData, &jsonUniCommand)

		if err != nil {
			log.Println("Invalid JSON data received")
			return
		}

		switch jsonUniCommand.Command {
		case UNI_COMMAND_REQUEST_ALL_DATA:
			g.sendConfigToClient(input.client)
			g.sendAllDataToClient(input.client)
		default:
			log.Println("Unknown command received")
		}

	default:
		log.Println("Unknown command type")
	}
}
