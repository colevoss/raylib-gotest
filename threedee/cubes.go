package main

import (
	"raylib-gotest/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_COLUMS = 25

type CubeSystem struct {
	heights   []float32
	positions []rl.Vector3
	colors    []rl.Color
}

func NewCubeSystem() *CubeSystem {
	return &CubeSystem{
		heights:   make([]float32, MAX_COLUMS),
		positions: make([]rl.Vector3, MAX_COLUMS),
		colors:    make([]rl.Color, MAX_COLUMS),
	}
}

func (cs *CubeSystem) Register(sm *engine.SystemManager) {
	sm.AddStartup(cs.Start)
	sm.AddCameraRender(cs.Draw)
}

func (cs *CubeSystem) Start() error {
	for i := 0; i < MAX_COLUMS; i++ {
		cs.heights[i] = float32(rl.GetRandomValue(1, 12))
		cs.positions[i] = rl.NewVector3(float32(rl.GetRandomValue(-15, 15)), cs.heights[i]/2, float32(rl.GetRandomValue(-15, 15)))
		cs.colors[i] = rl.NewColor(uint8(rl.GetRandomValue(20, 255)), uint8(rl.GetRandomValue(10, 55)), 30, 255)
	}

	mat := rl.LoadMaterialDefault()
	mmap := mat.GetMap(rl.MapDiffuse)
	mmap.Color = rl.Blue

	return nil
}

func (cs *CubeSystem) Draw() error {
	for i := 0; i < MAX_COLUMS; i++ {
		rl.DrawCube(cs.positions[i], 2.0, cs.heights[i], 2.0, cs.colors[i])
		rl.DrawCubeWires(cs.positions[i], 2.0, cs.heights[i], 2.0, rl.Blue)
	}
	return nil
}
