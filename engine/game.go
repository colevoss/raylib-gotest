package engine

import (
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Screen struct {
	Width  int32
	Height int32
}

type Game interface {
	ScreenWidth() int32
	ScreenHeight() int32
}

type GameImpl struct {
	sm        *SystemManager
	screen    *Screen
	name      string
	frameRate int32

	systems []System

	start    time.Time
	lastTick time.Time
}

func NewGame(width int32, height int32, name string) *GameImpl {
	screen := Screen{width, height}

	return &GameImpl{
		sm:        NewSystemManager(10),
		screen:    &screen,
		name:      name,
		frameRate: 60,
	}
}

func (g *GameImpl) ScreenWidth() int32 {
	return g.screen.Width
}

func (g *GameImpl) ScreenHeight() int32 {
	return g.screen.Height
}

func (g *GameImpl) Init() {
	rl.InitWindow(g.screen.Width, g.screen.Height, g.name)
	rl.SetTargetFPS(g.frameRate)
}

func (g *GameImpl) RegisterSystem(system System) {
	system.Register(g.sm)
}

func (g *GameImpl) Run() {
	g.start = time.Now()
	g.lastTick = time.Now()
	deltaTime := g.lastTick.Sub(g.start)

	g.sm.Pool.Start()
	g.sm.RunStartup()

	for !rl.WindowShouldClose() {
		// input

		g.sm.RunUpdate(deltaTime)
		g.sm.RunLateUpdate()

		// start drawing
		rl.BeginDrawing()
		rl.ClearBackground(rl.White)

		// g.systemManager.

		// engine playground

		g.sm.RunRender()
		// g.systemManager.Run(deltaTime)

		// let user camera stuff happen here
		// ...

		// end camera stuff somehow

		rl.DrawFPS(10, 10)
		rl.EndDrawing()

		g.sm.RunPostRender()

		deltaTime = time.Now().Sub(g.lastTick)
		g.lastTick = time.Now()
	}
}
