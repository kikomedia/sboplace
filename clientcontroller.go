package main

type ClientController struct {
	client          *Client
	gameController  *GameController
	stateController GameStateController
}

func (c *ClientController) changeState(state string) {

	if c.stateController != nil {
		c.stateController.onLeaveState()
	}

	switch state {
	case CLIENT_STATE_START:
		c.stateController = &ClientStateStart{ClientStateBase{controller: c.gameController, client: c.client}}

	default:
		return
	}
	c.client.data.state = state
	c.stateController.onEnterState()
}
