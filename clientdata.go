package main

import (
	"image/color"
	"math/rand"
	"time"
)

const (
	// Game Object Type
	GOT_PLAYER   = "player"
	GOT_GRAPHICS = "gfx"
)

type Position2D struct {
	x int64
	y int64
}

type ClientData struct {
	clientUUID    int64
	clientName    string
	position      Position2D
	state         string
	colorID       color.RGBA
	lastColorTime time.Time
	token         int64
}

func newClientData(clientUUID int64, clientId string) *ClientData {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return &ClientData{
		clientUUID: clientUUID,
		clientName: clientId,
		position:   Position2D{0, 0},
		colorID:    color.RGBA{R: 0, G: 0, B: 0, A: 0xFF},
		token:      r1.Int63n(99999999) + 1,
	}
}
