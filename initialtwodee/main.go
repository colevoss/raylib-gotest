package main

import (
	"fmt"
	"raylib-gotest/engine"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Counter struct {
	CountChan chan int
}

func (c *Counter) Add(i int) {
	c.CountChan <- i
}

type DisplaySystem struct {
	engine.System

	sync.RWMutex
	circleCount int
	msg         *Counter
}

func NewDisplaySystem(circleCount int, count *Counter) *DisplaySystem {
	return &DisplaySystem{
		circleCount: circleCount,
		msg:         count,
	}
}

func (d *DisplaySystem) Register(manager *engine.SystemManager) {
	manager.AddStartup(d.Start)
	manager.AddRender(d.Draw)
}

func (d *DisplaySystem) Start() error {
	go d.RunMsg()

	return nil
}

func (d *DisplaySystem) RunMsg() {
	for coundDelta := range d.msg.CountChan {
		d.Lock()

		d.circleCount += coundDelta

		d.Unlock()
	}
}

func (d *DisplaySystem) Draw() error {
	defer d.RUnlock()

	d.RLock()

	rl.DrawText(fmt.Sprintf("Alive: %d", d.circleCount), 25, 25, 30, rl.Black)

	return nil
}

const CIRCLES = 1000
const SYSTEMS = 25

func main() {
	game := engine.NewGame(800, 450, "Test Engine")

	game.Init()

	displayChannel := &Counter{
		CountChan: make(chan int),
	}

	displaySystem := NewDisplaySystem(CIRCLES*SYSTEMS, displayChannel)

	for range SYSTEMS {
		mySystem := NewMySystem(CIRCLES, game, displayChannel)
		game.RegisterSystem(mySystem)
	}

	game.RegisterSystem(displaySystem)

	game.Run()
}
