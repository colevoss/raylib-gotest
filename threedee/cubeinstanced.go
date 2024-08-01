package main

import (
	"raylib-gotest/engine"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const MAX_INSTANCES = 100000

type InstancedCubeSystem struct {
	transforms []rl.Matrix
	mesh       rl.Mesh
	material   rl.Material
	mmap       *rl.MaterialMap
	light      *Light
}

func NewInstancedCubeSystem() *InstancedCubeSystem {
	return &InstancedCubeSystem{
		transforms: make([]rl.Matrix, MAX_INSTANCES),
	}
}

func (s *InstancedCubeSystem) Register(sm *engine.SystemManager) {
	sm.AddStartup(s.start)
	sm.AddUpdate(s.update)
	sm.AddCameraRender(s.draw)
}

func (s *InstancedCubeSystem) start() error {
	rotations := make([]rl.Matrix, MAX_INSTANCES)    // Rotation state of instances
	rotationsInc := make([]rl.Matrix, MAX_INSTANCES) // Per-frame rotation animation of instances
	translations := make([]rl.Matrix, MAX_INSTANCES) // Locations of instances

	for i := 0; i < MAX_INSTANCES; i++ {
		x := float32(rl.GetRandomValue(-150, 150))
		y := float32(rl.GetRandomValue(-150, 150))
		z := float32(rl.GetRandomValue(-150, 150))
		translations[i] = rl.MatrixTranslate(x, y, z)

		x = float32(rl.GetRandomValue(0, 360))
		y = float32(rl.GetRandomValue(0, 360))
		z = float32(rl.GetRandomValue(0, 360))
		axis := rl.Vector3Normalize(rl.NewVector3(x, y, z))
		angle := float32(rl.GetRandomValue(0, 10)) * rl.Deg2rad

		rotationsInc[i] = rl.MatrixRotate(axis, angle)
		rotations[i] = rl.MatrixIdentity()
		rotations[i] = rl.MatrixMultiply(rotations[i], rotationsInc[i])
		s.transforms[i] = rl.MatrixMultiply(rotations[i], translations[i])
	}

	shader := rl.LoadShader("threedee/glsl330/base_lighting_instanced.vs", "threedee/glsl330/lighting.fs")
	shader.UpdateLocation(rl.ShaderLocMatrixMvp, rl.GetShaderLocation(shader, "mvp"))
	shader.UpdateLocation(rl.ShaderLocVectorView, rl.GetShaderLocation(shader, "viewPos"))
	shader.UpdateLocation(rl.ShaderLocMatrixModel, rl.GetShaderLocationAttrib(shader, "instanceTransform"))

	ambientLoc := rl.GetShaderLocation(shader, "ambient")
	rl.SetShaderValue(shader, ambientLoc, []float32{0.2, 0.2, 0.2, 1.0}, rl.ShaderUniformVec4)
	s.light = NewLight(LightTypeDirectional, rl.NewVector3(0.0, 150.0, 0.0), rl.Vector3Zero(), rl.White, shader)

	s.mesh = rl.GenMeshSphere(2.5, 10, 10)
	s.material = rl.LoadMaterialDefault()
	s.material.Shader = shader
	s.material.GetMap(rl.MapDiffuse).Color = rl.Blue
	s.material.GetMap(rl.MapMetalness).Value = 100.0

	return nil
}

func (s *InstancedCubeSystem) update(dt time.Duration) error {
	s.light.position.X += 10 * rl.GetFrameTime()
	return nil
}

func (s *InstancedCubeSystem) draw() error {
	rl.DrawMeshInstanced(s.mesh, s.material, s.transforms, MAX_INSTANCES)
	s.light.UpdateValues()
	return nil
}
