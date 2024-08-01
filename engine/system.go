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

type SystemType = int

const (
	Startup SystemType = iota
	Update
	LateUpdate
	Render
	PostRender
)

type SystemManager struct {
	startup    []SystemFunc
	update     []UpdateFunc
	lateUpdate []SystemFunc
	render     []SystemFunc
	postRender []SystemFunc
	Pool       *Pool
}

func NewSystemManager(poolCount uint) *SystemManager {
	return &SystemManager{
		Pool: NewPool(poolCount),
	}
}

func (sm *SystemManager) Startup(system SystemFunc) {
	sm.startup = append(sm.startup, system)
}

func (sm *SystemManager) Update(system UpdateFunc) {
	sm.update = append(sm.update, system)
}

func (sm *SystemManager) LateUpdate(system SystemFunc) {
	sm.lateUpdate = append(sm.lateUpdate, system)
}

func (sm *SystemManager) Render(system SystemFunc) {
	sm.render = append(sm.render, system)
}

func (sm *SystemManager) PostRender(system SystemFunc) {
	sm.postRender = append(sm.postRender, system)
}

func (sm *SystemManager) RunStartup() {
	for _, system := range sm.startup {
		system.Run()
	}
}

func (sm *SystemManager) RunUpdate(dt time.Duration) {
	if len(sm.update) == 0 {
		return
	}

	for _, system := range sm.update {
		sm.Pool.Add(func() error {
			system.Update(dt)
			return nil
		})
	}

	sm.Pool.Wait()
}

func (sm *SystemManager) RunLateUpdate() {
	if len(sm.lateUpdate) == 0 {
		return
	}

	for _, system := range sm.lateUpdate {
		sm.Pool.Add(system.Run)
	}

	sm.Pool.Wait()
}

func (sm *SystemManager) RunRender() {
	if len(sm.render) == 0 {
		return
	}

	for _, system := range sm.render {
		system.Run()
	}
}

func (sm *SystemManager) RunPostRender() {
	if len(sm.postRender) == 0 {
		return
	}

	for _, system := range sm.postRender {
		sm.Pool.Add(system.Run)
	}

	sm.Pool.Wait()
}

func (sm *SystemManager) Run(dt time.Duration) {
	sm.RunUpdate(dt)
	sm.RunLateUpdate()
	sm.RunRender()
	sm.RunPostRender()
}

func (sm *SystemManager) RunWithoutRender(dt time.Duration) {
	sm.RunUpdate(dt)
	sm.RunLateUpdate()
	// sm.RunRender()
	sm.RunPostRender()
}

type System interface {
	Register(sm *SystemManager)
}
