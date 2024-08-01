package main

import (
	"fmt"
	"math/rand"
	"raylib-gotest/engine"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func randomizeColor(color *rl.Color) {
	color.R = uint8(rl.GetRandomValue(55, 255))
	color.G = uint8(rl.GetRandomValue(55, 255))
	color.B = uint8(rl.GetRandomValue(55, 255))
}

type CircleCountSystem interface {
	Add(int)
}

type MySystem struct {
	engine.System

	*sync.RWMutex

	countChan    CircleCountSystem
	game         engine.Game
	maxHits      uint8
	CirlcesCount uint

	r *rand.Rand

	Pos []rl.Vector2

	alive *engine.SparseSet
	dead  *engine.SparseSet

	R []float32

	Dx []float32
	Dy []float32

	Colors []rl.Color

	HitCount []uint8
}

func NewMySystem(circles uint, game engine.Game, counter CircleCountSystem) *MySystem {
	r := rand.New(rand.NewSource(99))

	mySystem := &MySystem{
		countChan:    counter,
		RWMutex:      &sync.RWMutex{},
		maxHits:      2,
		game:         game,
		r:            r,
		CirlcesCount: circles,

		Pos: make([]rl.Vector2, circles),

		R: make([]float32, circles),

		Dx: make([]float32, circles),
		Dy: make([]float32, circles),

		alive: engine.NewSparseSet(circles + 1),
		dead:  engine.NewSparseSet(circles + 1),

		Colors: make([]rl.Color, circles),

		HitCount: make([]uint8, circles),
	}

	return mySystem
}

func (s *MySystem) Register(manager *engine.SystemManager) {
	manager.AddStartup(s.Start)
	manager.AddUpdate(s.Update)
	manager.AddRender(s.Draw)
	manager.AddPostRender(s.KillCircles)
}

func (s *MySystem) StartTimer() {
	fmt.Printf("Starting timer\n")
	timer := time.NewTicker(1 * time.Second)

	for {
		select {
		case <-timer.C:
			s.Lock()

			fmt.Printf("Attempting to spawn circle %d\n", s.alive.N())
			id := s.alive.N()

			if id < s.CirlcesCount && id < s.alive.Max() {
				fmt.Printf("Spawing circle %d\n", id)
				s.SpawnCircle(id)
				s.dead.Remove(id)
			}

			s.Unlock()
		}
	}
}

func (s *MySystem) KillCircles() error {
	for i := range s.alive.N() {
		id := s.alive.Get(i)

		if s.HitCount[id] > s.maxHits {
			s.KillCircle(id)
		}
	}

	return nil
}

func (s *MySystem) Start() error {
	height := s.game.ScreenHeight()
	width := s.game.ScreenWidth()

	for i := range s.CirlcesCount {
		id := i
		s.SpawnCircle(id)
		s.ensureCircleIsInBounds(id, height, width)
	}

	// go s.StartTimer()

	return nil
}

func (s *MySystem) SpawnCircle(id uint) {
	r := s.r.Float32() * 20
	s.R[id] = r

	vec := rl.Vector2{
		X: s.r.Float32()*20 + r,
		Y: s.r.Float32()*20 + r,
	}

	s.Pos[id] = vec

	s.Dx[id] = float32(rl.GetRandomValue(-10, 10)+10) * 10
	s.Dy[id] = float32(rl.GetRandomValue(-10, 10)+10) * 10
	s.Colors[id] = rl.Blue

	s.HitCount[id] = 0

	s.alive.AddGrow(id)
}

func (s *MySystem) KillCircle(id uint) {
	// fmt.Printf("Killing circle %d\n", id)
	s.alive.Remove(id)
	s.countChan.Add(-1)
}

func (s *MySystem) ensureCircleIsInBounds(id uint, height int32, width int32) {
	pos := s.Pos[id]
	r := s.R[id]

	if pos.X+r >= float32(width) {
		pos.X = float32(width)
	}

	if pos.X-r <= 0 {
		pos.X = 0
	}

	if pos.Y+r >= float32(height) {
		pos.Y = float32(height)
	}

	if pos.Y <= 0 {
		pos.Y = 0
	}
}

func (s *MySystem) MoveCircle(i uint, delta time.Duration) {
	deltaMs := float32(delta.Seconds())
	s.Pos[i].X += s.Dx[i] * deltaMs
	s.Pos[i].Y += s.Dy[i] * deltaMs
}

func (s *MySystem) CheckCircleBounds(i uint, height int32, width int32) {
	vec := s.Pos[i]
	r := s.R[i]

	if vec.X+r >= float32(width) || vec.X-r <= 0 {
		s.Dx[i] *= -1
		s.HitCount[i] += 1
		randomizeColor(&s.Colors[i])
	}

	if vec.Y+r >= float32(height) || vec.Y-r <= 0 {
		s.Dy[i] *= -1
		s.HitCount[i] += 1
		randomizeColor(&s.Colors[i])
	}
}

func (s *MySystem) Draw() error {
	for i := range s.alive.N() {
		id := s.alive.Get(i)
		s.DrawCircle(id)
	}

	return nil
}

func (s *MySystem) DrawCircle(i uint) {
	rl.DrawCircleV(
		s.Pos[i],
		s.R[i],
		s.Colors[i],
	)
}

func (s *MySystem) Update(deltaMs time.Duration) error {
	width := s.game.ScreenWidth()
	height := s.game.ScreenHeight()

	for i := range s.alive.N() {
		id := s.alive.Get(i)
		s.CheckCircleBounds(id, height, width)
		s.MoveCircle(id, deltaMs)
	}

	return nil
}
