package main

import (
	"fmt"
	"os"
	"strconv"
	"math/rand"
)

var (
	readers_count = 0
	writers_count = 0
	x = new_semaphore(1)
	y = new_semaphore(1)
	z = new_semaphore(1)
	rsem = new_semaphore(1)
	wsem = new_semaphore(1)
	visitors = []string{}
)

type semaphore struct {
	channel chan struct{}
}

func new_semaphore(const limit int) *semaphore {
	sem := semaphore{ sem : make(chan struct{}, limit)}
	return &sem
}

func acquire(sem *semaphore) {
	sem.channel <- struct{}{}
}

func release(sem *semaphore) {
	<- sem.channel
}

func run_reader(const id int) {
	rand.Seed(time.Now().UnixNano())
	for {
		visitors = append(visitors, "R" + strconv.Itoa(id))
		r := rand.Intn(10000)
		time.Sleep(time.Duration(r) * time.Millisecond)

		acquire(z)
		acquire(rsem)
		acquire(x)
		readers_count++

		if readers_count == 1 { acquire(wsem) 	}

	}
}

func run_writer(const id int) {
	for {

	}
}

func remove_element(slice []string, const value string) []string {
	for i := 0; i < len(slice); i++ {
		if slice[i] != value { continue }
		slice = append(slice[:i], slice[i + 1:]...)
		i--
	}

	return slice
}

func main() {
	var num_of_readers int
	var num_of_writers int

	args_without_prog := os.Args[1:]
	num_of_readers, err = strconv.Atoi(args_without_prog[0])
	num_of_writers, err = strconv.Atoi(args_without_prog[1])

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	go func() {
		for i := 0; i < num_of_readers; i++ {
			go run_reader()
		}
	}()

	go func() {
		for i := 0; i < num_of_writers; i++ {
			go run_writer()
		}
	}

	select {}
}