package main

import (
	"image/color"
	"raylib-gotest/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GroundSystem struct {
	center rl.Vector3
	size   rl.Vector2
	color  color.RGBA

	walls [3]Wall
}

type Wall struct {
	pos    rl.Vector3
	height float32
	width  float32
	length float32
	color  color.RGBA
}

func NewGrounSystem() *GroundSystem {
	system := &GroundSystem{
		center: rl.NewVector3(0.0, 0.0, 0.0),
		size:   rl.NewVector2(32.0, 32.0),
		color:  rl.LightGray,
	}

	system.walls[0] = Wall{
		pos:    rl.NewVector3(-16.0, 2.5, 0.0),
		width:  1.0,
		height: 5.0,
		length: 32.0,
		color:  rl.Blue,
	}

	system.walls[1] = Wall{
		pos:    rl.NewVector3(16.0, 2.5, 0.0),
		width:  1.0,
		height: 5.0,
		length: 32.0,
		color:  rl.Lime,
	}

	system.walls[2] = Wall{
		pos:    rl.NewVector3(0, 2.5, 16.0),
		width:  32.0,
		height: 5.0,
		length: 1.0,
		color:  rl.Gold,
	}

	return system
}

func (gs *GroundSystem) Register(sm *engine.SystemManager) {
	sm.AddCameraRender(gs.Draw)
}

func (gs *GroundSystem) Draw() error {
	rl.DrawPlane(gs.center, gs.size, gs.color)

	for _, wall := range gs.walls {
		rl.DrawCube(wall.pos, wall.width, wall.height, wall.length, wall.color)
	}

	return nil
}
