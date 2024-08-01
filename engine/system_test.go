package engine

import (
	"fmt"
	"testing"
	"time"
)

type Vec struct {
	X float32
	Y float32
}

type TestSystem struct {
	System
	Count uint
	Vecs  []Vec
	Alive *SparseSet
}

func NewTestSystem(count uint) *TestSystem {
	return &TestSystem{
		Count: count,
		Vecs:  make([]Vec, count),
		Alive: NewSparseSet(count),
	}
}

func (t *TestSystem) Register(sm *SystemManager) {
	sm.Startup(t.Start)
	sm.Update(t.Update)
	sm.PostRender(t.PostRender)
}

func (t *TestSystem) Start() error {
	for i := range t.Count {
		id := i

		t.Vecs = append(t.Vecs, Vec{
			X: float32(id),
			Y: float32(id),
		})

		t.Alive.Add(id)
	}

	return nil
}

func (t *TestSystem) Update(dt time.Duration) error {
	for i := range t.Alive.N() {
		id := t.Alive.Get(i)

		t.Vecs[id].X /= float32(i)
		t.Vecs[id].Y /= float32(i)
	}

	return nil
}

func (t *TestSystem) PostRender() error {
	for i := range t.Alive.N() {
		id := t.Alive.Get(i)

		if t.Vecs[id].X < 0.10 {
			// t.Alive.Remove(i)
			t.Vecs[id].X += 100
			t.Vecs[id].Y += 100
		}
	}

	return nil
}

func benchmarkSystem(vecCount uint, systemCount uint, runCount uint, poolCount uint, b *testing.B) {
	dt, _ := time.ParseDuration("1ms")

	manager := NewSystemManager(poolCount)

	for range systemCount {
		s := NewTestSystem(vecCount)
		s.Register(manager)
	}

	manager.Pool.Start()
	manager.RunStartup()

	for n := 0; n < b.N; n++ {

		// for range runCount {
		manager.Run(dt)
		// }
	}
}

var Runs = []struct {
	vecs    int
	systems int
	pools   int
}{
	{vecs: 10, systems: 1, pools: 1},
	{vecs: 100, systems: 1, pools: 1},
	{vecs: 1000, systems: 1, pools: 1},
	{vecs: 10000, systems: 1, pools: 1},
	{vecs: 50000, systems: 1, pools: 1},

	{vecs: 10, systems: 1, pools: 10},
	{vecs: 100, systems: 1, pools: 10},
	{vecs: 1000, systems: 1, pools: 10},
	{vecs: 10000, systems: 1, pools: 10},
	{vecs: 50000, systems: 1, pools: 10},

	{vecs: 10, systems: 10, pools: 1},
	{vecs: 100, systems: 10, pools: 1},
	{vecs: 1000, systems: 10, pools: 1},
	{vecs: 10000, systems: 10, pools: 1},
	{vecs: 50000, systems: 10, pools: 1},

	{vecs: 10, systems: 10, pools: 10},
	{vecs: 100, systems: 10, pools: 10},
	{vecs: 1000, systems: 10, pools: 10},
	{vecs: 10000, systems: 10, pools: 10},
	{vecs: 50000, systems: 10, pools: 10},

	{vecs: 10, systems: 100, pools: 1},
	{vecs: 100, systems: 100, pools: 1},
	{vecs: 1000, systems: 100, pools: 1},
	{vecs: 10000, systems: 100, pools: 1},
	{vecs: 50000, systems: 100, pools: 1},

	{vecs: 10, systems: 100, pools: 10},
	{vecs: 100, systems: 100, pools: 10},
	{vecs: 1000, systems: 100, pools: 10},
	{vecs: 10000, systems: 100, pools: 10},
	{vecs: 50000, systems: 100, pools: 10},

	{vecs: 10, systems: 10, pools: 100},
	{vecs: 100, systems: 10, pools: 100},
	{vecs: 1000, systems: 10, pools: 100},
	{vecs: 10000, systems: 10, pools: 100},
	{vecs: 50000, systems: 10, pools: 100},

	{vecs: 10, systems: 100, pools: 100},
	{vecs: 100, systems: 100, pools: 100},
	{vecs: 1000, systems: 100, pools: 100},
	{vecs: 10000, systems: 100, pools: 100},
	{vecs: 50000, systems: 100, pools: 100},
}

func BenchmarkSystemsWithPool(b *testing.B) {
	for _, run := range Runs {
		b.Run(fmt.Sprintf("%d_Vecs__%d_Systems__%d_Pool", run.vecs, run.systems, run.pools), func(b *testing.B) {
			manager := NewSystemManager(uint(run.pools))

			for range run.systems {
				s := NewTestSystem(uint(run.vecs))
				s.Register(manager)
			}

			manager.Pool.Start()
			manager.RunStartup()

			dt, _ := time.ParseDuration("1ms")
			for n := 0; n < b.N; n++ {
				manager.Run(dt)
			}
		})
	}
}
