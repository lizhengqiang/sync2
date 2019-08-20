package sync2

import "sync"

type Handler func(id int, task interface{}) (err error, retry bool)

type Pool struct {
	poolSize   int
	tasks      chan interface{}
	worker     Handler
	wg         *sync.WaitGroup
	errorStack []error
}

func NewPool(poolSize int, worker Handler) *Pool {
	return &Pool{
		poolSize:   poolSize,
		tasks:      make(chan interface{}, poolSize),
		worker:     worker,
		wg:         &sync.WaitGroup{},
		errorStack: []error{},
	}
}

func (s *Pool) Put(task interface{}) {
	s.wg.Add(1)
	s.tasks <- task
}

func (s *Pool) Start() {
	for i := 0; i < s.poolSize; i++ {
		workerId := i
		go func() {
			for {
				task := <-s.tasks
				err, retry := s.worker(workerId, task)
				if err != nil {
					s.errorStack = append(s.errorStack, err)
				}
				if retry {
					s.Put(task)
				}
				s.wg.Done()
			}
		}()
	}

}

func (s *Pool) Wait() {
	s.wg.Wait()
}
