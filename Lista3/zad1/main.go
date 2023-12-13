package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

var output = []string{}

const (
	THINKING = 0
	HUNGRY = 1
	EATING = 2
)

func remove_element(slice []string, value string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] != value { continue }
		slice = append(slice[:i], slice[i + 1:]...)
		i--
	}
	return slice
}

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

type cond struct {
	me semaphore
	waiters int
	channel chan struct{}
}

func new_cond() *cond {
	cond := cond{ channel : make(chan struct{}), me : *new_semaphore(1)}
	return &cond
}

func (c *cond) wait() {
	acquire(&c.me)
	c.waiters++
	localCh := c.channel
	release(&c.me)

	<-localCh

	acquire(&c.me)
	c.waiters--
	if c.waiters == 0 {
		c.channel = make(chan struct{})
	}
	release(&c.me)
}

func (c *cond) signal() {
	acquire(&c.me)
	defer release(&c.me)

	if c.waiters > 0 {
		close(c.channel)
	}
}

type dining_philosophers struct {
	state []int
	cond []*cond
	num_phils int
}

func new_dining_philosophers(num_phils int) *dining_philosophers {
	dp := &dining_philosophers{
		state:    make([]int, num_phils),
		cond:     make([]*cond, num_phils),
		num_phils: num_phils,
	}

	for i := 0; i < num_phils; i++ {
		dp.cond[i] = new_cond()
	}

	return dp
}

func left(id int) int {
	return id
}

func right(id int, max int) int {
	return (id + 1) % max
}

func (dp *dining_philosophers) take_forks(id int) {
	dp.state[id] = HUNGRY
	dp.test(id)

	if dp.state[id] != EATING { dp.cond[id].wait() }

	output = append(output, "W" +strconv.Itoa(left(id)) + " F" + strconv.Itoa(id) + " W" + strconv.Itoa(right(id, dp.num_phils)))
	fmt.Println(output)
}

func (dp *dining_philosophers) drop_forks(id int) {
	dp.state[id] = THINKING
	output = remove_element(output, "W" + strconv.Itoa(left(id)) + " F" + strconv.Itoa(id) + " W" + strconv.Itoa(right(id, dp.num_phils)))

	fmt.Println(output)

	dp.test((id + 1) % dp.num_phils)
	dp.test((id + dp.num_phils - 1) % dp.num_phils)
}

func (dp *dining_philosophers) test(id int) {
	if 	dp.state[(id + 1) % dp.num_phils] != EATING &&
		dp.state[(id + dp.num_phils - 1) % dp.num_phils] != EATING &&
		dp.state[id] == HUNGRY {
	
		dp.state[id] = EATING
		dp.cond[id].signal()
	}
}

func (dp *dining_philosophers) philosopher_life(id int) {
	rand.Seed(time.Now().UnixNano())
	for {	
		r := rand.Intn(1000)
		time.Sleep(time.Duration(r) * time.Millisecond)

		dp.take_forks(id)

		rand.Seed(time.Now().UnixNano())
		r = rand.Intn(1000)
		time.Sleep(time.Duration(r) * time.Millisecond)

		dp.drop_forks(id)
	}
}


func main() {
	num_philosophers := 5

	dining_philosophers := new_dining_philosophers(num_philosophers)

	for i := 0; i < num_philosophers; i++ {
		go dining_philosophers.philosopher_life(i)
	}

	select {}
}