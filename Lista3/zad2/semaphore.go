package main

type semaphore struct {
	channel chan struct{}
}

func new_semaphore(limit int) *semaphore {
	sem := semaphore{ channel : make(chan struct{}, limit)}
	return &sem
}

func acquire(sem *semaphore) {
	sem.channel <- struct{}{}
}

func release(sem *semaphore) {
	<- sem.channel
}