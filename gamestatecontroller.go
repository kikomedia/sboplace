package main

import (
	"log"
	"time"
)

const (
	GAME_STATE_START = "start"
	GAME_STATE_MAIN  = "main"
)

type GameStateController interface {
	onEnterState()
	onLeaveState()
	onUpdate(delta time.Duration)
	onProcessInput(sender *Client, commands []string)
	onProcessCommand(sender *Client, command SimpleCommandT, parameters []string)
}

type GameStateBase struct {
	controller *GameController
}

type GameStateMain struct {
	GameStateBase
}

type GameStateStart struct {
	GameStateBase
}

func (g *GameStateBase) onEnterState() {

}

func (g *GameStateBase) onLeaveState() {

}

func (g *GameStateBase) onUpdate(delta time.Duration) {

}

func (g *GameStateBase) onProcessInput(sender *Client, commands []string) {
	sender.clientStateController.onProcessInput(sender, commands)
}

func (g *GameStateBase) onProcessCommand(sender *Client, command SimpleCommandT, parameters []string) {
	sender.clientStateController.onProcessCommand(sender, command, parameters)
}

func (g *GameStateMain) onUpdate(delta time.Duration) {

}

func (g *GameStateStart) onUpdate(delta time.Duration) {

}

func (g *GameStateStart) onProcessInput(sender *Client, commands []string) {

	if len(commands) > 0 {

		switch commands[0] {
		case KEY_COMMAND_SPACE:
			g.controller.changeState(GAME_STATE_MAIN)
		default:
			log.Println("Unknown Command")
		}

	} else {
		log.Println("No parameters received")
	}
}
