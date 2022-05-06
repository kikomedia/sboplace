package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type GameState int

const STEP_SIZE = math.Pi / 45
const CLIENT_STEP_SIZE int64 = 8
const MAP_WIDTH = 100
const MAP_HEIGHT = 75
const MAX_COLORS = 16
const TIMELINE_FOLDER = "timeline"

const (
	GameStateMenu GameState = iota
	GameStateRunning
	GameStateFinished
)

type PlayerState int

const (
	PlayerStateFirstConnection PlayerState = iota
	PlayerStateReady
	PlayerRunning
)

type BlockData struct {
	blockUUID int64
	position  Position2D
	blockType string
	colorID   int64
}

type GameBlock struct {
	blockType string
	colorID   color.RGBA
}

type GameData struct {
	state            GameState
	stepper          float64
	blocks           []BlockData
	matrix           [][]GameBlock
	currentBlockUUID int64
}

func newGameData() *GameData {
	gd := &GameData{state: GameStateRunning, stepper: 0.0, blocks: []BlockData{}, currentBlockUUID: 0}

	gd.matrix = make([][]GameBlock, MAP_HEIGHT)

	for i := 0; i < MAP_HEIGHT; i++ {
		gd.matrix[i] = make([]GameBlock, MAP_WIDTH)
		for j := 0; j < MAP_WIDTH; j++ {
			gd.matrix[i][j].colorID = color.RGBA{R: 0, G: 0, B: 0, A: 0}
		}
	}

	return gd
}

func (b *GameData) getGridPositionByPlayerPosition(x int64, y int64) (int64, int64) {
	rx := int64(math.Round(float64(x / CLIENT_STEP_SIZE)))
	ry := int64(math.Round(float64(y / CLIENT_STEP_SIZE)))
	return rx, ry
}

func (b *GameData) addBlockOnMatrix(x int64, y int64, colorID color.RGBA) {
	//if b.isValidGridPos(x, y) {
	b.matrix[y][x].colorID = colorID
	//}
}

func (b *GameData) isValidGridPos(x int64, y int64) bool {
	if (x >= 0 && x < MAP_WIDTH) && (y >= 0 && y < MAP_HEIGHT) {
		return true
	}
	return false
}

func (b *GameData) getNewestFileName() (string, error) {
	dir := TIMELINE_FOLDER + "/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	var newestFile string
	var newestTime int64 = 0
	if len(files) == 0 {
		return "", nil
	}
	for _, f := range files {
		fi, err := os.Stat(dir + f.Name())
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		currTime := fi.ModTime().Unix()
		if currTime > newestTime {
			newestTime = currTime
			newestFile = f.Name()
		}
	}

	return newestFile, nil
}

func (b *GameData) loadFromFile(filename string) {
	infile, err := os.Open("timeline/" + filename)
	if err != nil {
		panic("Error on file")
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if err != nil {
		panic("Error on decode")
	}

	for y := 0; y < MAP_HEIGHT; y++ {
		for x := 0; x < MAP_WIDTH; x++ {
			//col := b.matrix[y][x].colorID
			col := src.At(x, y)

			cr, cg, cb, ca := col.RGBA()

			b.matrix[y][x].colorID.R = uint8(cr)
			b.matrix[y][x].colorID.G = uint8(cg)
			b.matrix[y][x].colorID.B = uint8(cb)
			b.matrix[y][x].colorID.A = uint8(ca)
		}
	}
}

func (b *GameData) saveToFile() {
	fmt.Println("Save to file...")
	img := image.NewRGBA(image.Rect(0, 0, MAP_WIDTH, MAP_HEIGHT))

	for y := 0; y < MAP_HEIGHT; y++ {
		for x := 0; x < MAP_WIDTH; x++ {
			col := b.matrix[y][x].colorID
			img.Set(x, y, color.RGBA{
				R: col.R,
				G: col.G,
				B: col.B,
				A: col.A,
			})
		}
	}

	ts := time.Now().UTC().Format(time.RFC3339)
	filename := strings.Replace(strings.Replace(ts, ":", "", -1), "-", "", -1)
	filename = filename + ".png"

	if _, err := os.Stat(TIMELINE_FOLDER); os.IsNotExist(err) {
		if err := os.Mkdir(TIMELINE_FOLDER, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.Create(TIMELINE_FOLDER + "/" + filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
