package main

import (
	"raylib-gotest/engine"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type CameraSystem struct {
	camera *rl.Camera3D
}

func NewCameraSystem(camera *rl.Camera3D) *CameraSystem {
	return &CameraSystem{
		camera: camera,
	}
}

func (cs *CameraSystem) Register(sm *engine.SystemManager) {
	sm.AddStartup(cs.Start)
	// sm.AddUpdate(cs.Update)
	sm.AddLateUpdate(cs.LateUpdate)
}

func (cs *CameraSystem) Start() error {
	cs.camera.Position = rl.NewVector3(0.0, 0.0, 4.0)
	cs.camera.Target = rl.NewVector3(0.0, 1.8, 0.0)
	cs.camera.Up = rl.NewVector3(0.0, 1.0, 0.0)
	cs.camera.Fovy = 45.0
	cs.camera.Projection = rl.CameraPerspective

	return nil
}

func (cs *CameraSystem) Update(dt time.Duration) error {
	if rl.IsKeyDown(rl.KeyW) {
		cs.camera.Position.X += 10 * float32(dt.Seconds())
	}

	return nil
}

func (cs *CameraSystem) LateUpdate() error {
	rl.UpdateCamera(cs.camera, rl.CameraFree)

	return nil
}
