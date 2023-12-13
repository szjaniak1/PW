package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"math/rand"
)

var (
	readers_count = 0
	writers_count = 0
	sh_var = 5
	rsem = new_semaphore(1)
	wsem = new_semaphore(1)
	visitors = []string{}
)

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

func run_reader(id int) {
	rand.Seed(time.Now().UnixNano())
	for {
		visitors = append(visitors, "R" + strconv.Itoa(id))
		fmt.Println(visitors)
		r := rand.Intn(10000)
		time.Sleep(time.Duration(r) * time.Millisecond)

		acquire(rsem)
		readers_count++

		if readers_count == 1 { acquire(wsem) 	}
		release(rsem)

		fmt.Println("updated val " + strconv.Itoa(sh_var))
		acquire(rsem)

		readers_count--
		if readers_count == 0 { release(wsem) }
		release(rsem)

		r = rand.Intn(10000)
		time.Sleep(time.Duration(r) * time.Millisecond)
	}
}

func run_writer(id int) {  
	rand.Seed(time.Now().UnixNano())
	for {
		visitors = append(visitors, "W" + strconv.Itoa(id))
		fmt.Println(visitors)
		r := rand.Intn(10000)
		time.Sleep(time.Duration(r) * time.Millisecond)

		writers_count++
		acquire(wsem)

		sh_var += 5
		release(wsem)
		writers_count--

		visitors = remove_element(visitors, "W" + strconv.Itoa(id))
		fmt.Println(visitors)
		r = rand.Intn(10000)
		time.Sleep(time.Duration(r) * time.Millisecond)
	}
}

func remove_element(slice []string, value string) []string {
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
	var err error

	args_without_prog := os.Args[1:]
	num_of_readers, err = strconv.Atoi(args_without_prog[0])
	num_of_writers, err = strconv.Atoi(args_without_prog[1])

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	go func() {
		for i := 0; i < num_of_readers; i++ {
			go run_reader(i)
		}
	}()

	go func() {
		for i := 0; i < num_of_writers; i++ {
			go run_writer(i)
		}
	}()

	select {}
}