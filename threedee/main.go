package main

import (
	"raylib-gotest/engine"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	systemManager := engine.NewSystemManager(1)

	var screenWidth int32 = 800
	var screenHeight int32 = 450

	rl.InitWindow(screenWidth, screenHeight, "Test Three Dee")
	rl.SetTargetFPS(60)
	rl.SetConfigFlags(rl.FlagMsaa4xHint) // Enable Multi Sampling Anti Aliasing 4x (if available)

	camera := rl.Camera3D{}

	cameraSystem := NewCameraSystem(&camera)
	cameraSystem.Register(systemManager)

	// groundSystem := NewGrounSystem()
	// groundSystem.Register(systemManager)

	instancedCubeSystem := NewInstancedCubeSystem()
	instancedCubeSystem.Register(systemManager)

	// cubeSystem := NewCubeSystem()
	// cubeSystem.Register(systemManager)

	start := time.Now()
	lastTick := time.Now()
	deltaTime := lastTick.Sub(start)

	systemManager.Start()
	systemManager.Startup()

	for !rl.WindowShouldClose() {
		systemManager.Update(deltaTime)
		systemManager.LateUpdate()

		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.RayWhite)

			rl.BeginMode3D(camera)
			{
				systemManager.CameraRender()
			}

			rl.EndMode3D()

			systemManager.Render()

			rl.DrawFPS(10, 10)
		}
		rl.EndDrawing()
	}
}
