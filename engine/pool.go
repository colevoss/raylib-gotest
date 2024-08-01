package engine

import (
	"sync"
)

type Task func() error

type Job func() error

func (j Job) Run(wg *sync.WaitGroup) error {
	defer wg.Done()

	return j()
}

type Pool struct {
	wg      sync.WaitGroup
	queue   chan Job
	workers uint
}

func NewPool(workers uint) *Pool {
	return &Pool{
		workers: workers,
		queue:   make(chan Job),
	}
}

func (p *Pool) Add(job Job) {
	p.wg.Add(1)
	p.queue <- job
}

func (p *Pool) Worker(i int) {
	for job := range p.queue {
		job.Run(&p.wg)
	}
}

func (p *Pool) Start() {
	i := 1
	for range p.workers {
		go p.Worker(i)
		i++
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
