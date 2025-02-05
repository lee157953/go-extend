package ext

import (
	"sync"
	"sync/atomic"
)

const (
	_PromisePending   = 0
	_PromiseCompleted = 1
	_PromiseCanceled  = -1
)

type Promise[T any] struct {
	*_PromisePin[T]
}

type _PromisePin[T any] struct {
	locker sync.Mutex
	waiter sync.WaitGroup
	status int8
	result T
}

func (p Promise[T]) Pending() bool {
	p.locker.Lock()
	defer p.locker.Unlock()
	return p.status == _PromisePending
}

func (p Promise[T]) Canceled() bool {
	p.locker.Lock()
	defer p.locker.Unlock()
	return p.status == _PromiseCanceled
}

func (p Promise[T]) Completed() bool {
	p.locker.Lock()
	defer p.locker.Unlock()
	return p.status == _PromiseCompleted
}

func (p Promise[T]) Done() bool {
	return !p.Pending()
}

func (p Promise[T]) Cancel() bool {
	if p.status == _PromisePending {
		p.locker.Lock()
		defer p.waiter.Done()
		defer p.locker.Unlock()
		if p.status == _PromisePending {
			p.status = _PromiseCanceled
			return true
		}
	}
	return false
}

func (p Promise[T]) Complete(t T) bool {
	if p.status == _PromisePending {
		p.locker.Lock()
		defer p.waiter.Done()
		defer p.locker.Unlock()
		if p.status == _PromisePending {
			p.result = t
			p.status = _PromiseCompleted
			return true
		}
	}
	return false
}

func (p Promise[T]) Await() Opt[T] {
	p.waiter.Wait()
	if p.status == _PromiseCompleted {
		return Some(p.result)
	}
	return None[T]()
}

func (p Promise[T]) TryGet() Opt[T] {
	p.locker.Lock()
	defer p.locker.Unlock()
	if p.status == _PromiseCompleted {
		return Some(p.result)
	}
	return None[T]()
}

func Promise_[T any]() Promise[T] {
	p := Promise[T]{&_PromisePin[T]{}}
	p.waiter.Add(1)
	return p
}

type Canceler struct {
	b *atomic.Bool
}

func (c Canceler) Cancel() {
	c.b.Store(true)
}

func (c Canceler) Canceled() bool {
	return c.b.Load()
}

func Canceller_() Canceler {
	c := Canceler{new(atomic.Bool)}
	return c
}
