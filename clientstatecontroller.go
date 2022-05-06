package main

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
	"time"
)

const (
	CLIENT_STATE_START = "start"
	CLIENT_STATE_MAIN  = "main"
)

type ClientStateBase struct {
	controller *GameController
	client     *Client
}

type ClientStateStart struct {
	ClientStateBase
}

func (g *ClientStateStart) onEnterState() {

}

func (g *ClientStateStart) onLeaveState() {

}

func (g *ClientStateStart) onUpdate(delta time.Duration) {

}

func (g *ClientStateStart) onProcessCommand(sender *Client, command SimpleCommandT, parameters []string) {

	if command == SIMPLE_COMMAND_SET_COLOR {
		if len(parameters) == 1 {

			sender.data.colorID = ParseStringToColor(parameters[0])
			log.Println("Set Color:", sender.data.colorID)
			g.controller.server.sendToAllClients(createClientDataUpdateCommand(sender))
		}
	}
}

func (g *ClientStateStart) onProcessInput(sender *Client, commands []string) {
	if len(commands) > 0 {
		for _, command := range commands {
			switch command {
			case KEY_COMMAND_ESCAPE:
				log.Println("No action")
				//g.controller.changeState("start")
				//message := CreateUniCommand(UNI_COMMAND_CLEAR_ALL_DATA)
				//g.controller.server.sendToAllClients(message)
			case KEY_COMMAND_UP:
				sender.data.position.y -= CLIENT_STEP_SIZE
				if sender.data.position.y < 0 {
					sender.data.position.y = 0
				}
				g.controller.server.sendToAllClients(createClientDataUpdateCommand(sender))
				fmt.Println(sender.data.position)
			case KEY_COMMAND_DOWN:
				sender.data.position.y += CLIENT_STEP_SIZE
				if sender.data.position.y > MAP_HEIGHT*CLIENT_STEP_SIZE-CLIENT_STEP_SIZE {
					sender.data.position.y = MAP_HEIGHT*CLIENT_STEP_SIZE - CLIENT_STEP_SIZE
				}
				g.controller.server.sendToAllClients(createClientDataUpdateCommand(sender))
				fmt.Println(sender.data.position)
			case KEY_COMMAND_LEFT:
				sender.data.position.x -= CLIENT_STEP_SIZE
				if sender.data.position.x < 0 {
					sender.data.position.x = 0
				}
				g.controller.server.sendToAllClients(createClientDataUpdateCommand(sender))
				fmt.Println(sender.data.position)
			case KEY_COMMAND_RIGHT:
				sender.data.position.x += CLIENT_STEP_SIZE
				if sender.data.position.x > MAP_WIDTH*CLIENT_STEP_SIZE-CLIENT_STEP_SIZE {
					sender.data.position.x = MAP_WIDTH*CLIENT_STEP_SIZE - CLIENT_STEP_SIZE
				}
				g.controller.server.sendToAllClients(createClientDataUpdateCommand(sender))
				fmt.Println(sender.data.position)
			case KEY_COMMAND_RETURN:
				log.Println("New colorID :", sender.data.colorID)
				g.controller.server.sendToAllClients(createClientDataUpdateCommand(sender))
			case KEY_COMMAND_BACK:
				x, y := g.controller.server.gamedata.getGridPositionByPlayerPosition(sender.data.position.x, sender.data.position.y)

				diff := time.Now().Sub(sender.data.lastColorTime)

				if diff > time.Millisecond*50 {
					sender.data.lastColorTime = time.Now()

					g.controller.server.gamedata.matrix[y][x] = GameBlock{}
					g.controller.server.gamedata.matrix[y][x].colorID = color.RGBA{R: 0, G: 0, B: 0, A: 0}

					message := CreateSimpleCommand(SIMPLE_COMMAND_REMOVE_OBJECT_DATA, y*MAP_WIDTH+x)
					log.Println("Remove on :", x, y)
					g.controller.server.sendToAllClients(message)
				}

			case KEY_COMMAND_SPACE:
				x, y := g.controller.server.gamedata.getGridPositionByPlayerPosition(sender.data.position.x, sender.data.position.y)

				diff := time.Now().Sub(sender.data.lastColorTime)

				if diff > time.Millisecond*50 {
					sender.data.lastColorTime = time.Now()
					g.controller.server.gamedata.addBlockOnMatrix(x, y, sender.data.colorID)

					message := CreateSimpleCommand(SIMPLE_COMMAND_ADD_OBJECT_DATA, y*MAP_WIDTH+x, "block",
						strconv.FormatInt(x*CLIENT_STEP_SIZE, 10),
						strconv.FormatInt(y*CLIENT_STEP_SIZE, 10),
						ParseColorToString(sender.data.colorID),
					)

					g.controller.server.sendToAllClients(message)
				}

			default:
				log.Println("Unknown Command: ", command)
			}
		}

	} else {
		log.Println("No parameters received")
	}
}
