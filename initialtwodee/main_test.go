package main

import (
	"fmt"
	"raylib-gotest/engine"
	"testing"
	"time"
)

type MockGame struct {
}

func (mg *MockGame) ScreenWidth() int32 {
	return 100
}

func (mg *MockGame) ScreenHeight() int32 {
	return 100
}

var systemCases = []struct{ circles uint }{
	{circles: 10},
	{circles: 100},
	{circles: 1000},
	{circles: 2000},
	{circles: 5000},
	{circles: 10000},
	{circles: 20000},
	{circles: 50000},
	{circles: 100000},
	{circles: 200000},
	{circles: 500000},
}

type MockCountSystem struct {
}

func (m *MockCountSystem) Add(i int) {

}

func BenchmarkMySystem(b *testing.B) {
	for _, test := range systemCases {
		b.Run(fmt.Sprintf("circles_%d", test.circles), func(b *testing.B) {
			game := &MockGame{}
			// for n := 0; n < b.N; n++ {
			systemManager := engine.NewSystemManager(1)

			mySystem := NewMySystem(test.circles, game, &MockCountSystem{})
			mySystem.Register(systemManager)

			systemManager.Pool.Start()
			systemManager.RunStartup()

			dt, _ := time.ParseDuration("16ms")

			// for range 60 {
			for n := 0; n < b.N; n++ {
				systemManager.RunWithoutRender(dt)
			}
			// }
		})
	}
}
