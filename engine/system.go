package engine

import (
	"time"
)

type UpdateFunc func(dt time.Duration) error

func (uf UpdateFunc) Update(dt time.Duration) error {
	return uf(dt)
}

type SystemFunc func() error

func (sf SystemFunc) Run() error {
	return sf()
}

type SystemManager struct {
	startup      []SystemFunc
	update       []UpdateFunc
	lateUpdate   []SystemFunc
	render       []SystemFunc
	cameraRender []SystemFunc
	postRender   []SystemFunc
	Pool         *Pool
}

func NewSystemManager(poolCount uint) *SystemManager {
	return &SystemManager{
		Pool: NewPool(poolCount),
	}
}

func (sm *SystemManager) AddStartup(system SystemFunc) {
	sm.startup = append(sm.startup, system)
}

func (sm *SystemManager) AddUpdate(system UpdateFunc) {
	sm.update = append(sm.update, system)
}

func (sm *SystemManager) AddLateUpdate(system SystemFunc) {
	sm.lateUpdate = append(sm.lateUpdate, system)
}

func (sm *SystemManager) AddCameraRender(system SystemFunc) {
	sm.cameraRender = append(sm.cameraRender, system)
}

func (sm *SystemManager) AddRender(system SystemFunc) {
	sm.render = append(sm.render, system)
}

func (sm *SystemManager) AddPostRender(system SystemFunc) {
	sm.postRender = append(sm.postRender, system)
}

func (sm *SystemManager) Start() {
	sm.Pool.Start()
}

func (sm *SystemManager) Startup() {
	for _, system := range sm.startup {
		system.Run()
	}
}

func (sm *SystemManager) Update(dt time.Duration) {
	for _, system := range sm.update {
		sm.Pool.Add(func() error {
			system.Update(dt)
			return nil
		})
	}

	sm.Pool.Wait()
}

func (sm *SystemManager) LateUpdate() {
	for _, system := range sm.lateUpdate {
		sm.Pool.Add(system.Run)
	}

	sm.Pool.Wait()
}

func (sm *SystemManager) CameraRender() {
	for _, system := range sm.cameraRender {
		system.Run()
	}
}

func (sm *SystemManager) Render() {
	for _, system := range sm.render {
		system.Run()
	}
}

func (sm *SystemManager) PostRender() {
	for _, system := range sm.postRender {
		sm.Pool.Add(system.Run)
	}

	sm.Pool.Wait()
}

func (sm *SystemManager) Run(dt time.Duration) {
	sm.Update(dt)
	sm.LateUpdate()
	sm.CameraRender()
	sm.Render()
	sm.PostRender()
}

func (sm *SystemManager) RunWithoutRender(dt time.Duration) {
	sm.Update(dt)
	sm.LateUpdate()
	sm.PostRender()
}

type System interface {
	Register(sm *SystemManager)
}
