package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"math/rand"
)

const max_time_to_sleep = 10000

var (
	readers_count = 0
	book = 0
	rsem = new_semaphore(1)
	wsem = new_semaphore(1)
	visitors = []string{}
)

func run_reader(id int) {
	rand.Seed(time.Now().UnixNano())
	id_str := "R" + strconv.Itoa(id)

	for {
		r := rand.Intn(max_time_to_sleep)
		time.Sleep(time.Duration(r) * time.Millisecond)

		acquire(rsem)

		readers_count++
		if readers_count == 1 { acquire(wsem) }
		visitors = append(visitors, id_str)
		fmt.Println(visitors)
		fmt.Println("read book: " + strconv.Itoa(book))

		release(rsem)

		r = rand.Intn(max_time_to_sleep)
		time.Sleep(time.Duration(r) * time.Millisecond)

		acquire(rsem)

		readers_count--
		if readers_count == 0 { release(wsem) }
		visitors = remove_element(visitors, id_str)
		fmt.Println(visitors)

		release(rsem)
	}
}

func run_writer(id int) {  
	rand.Seed(time.Now().UnixNano())
	id_str := "W" + strconv.Itoa(id)

	for {
		r := rand.Intn(max_time_to_sleep)
		time.Sleep(time.Duration(r) * time.Millisecond)

		acquire(wsem)
		visitors = append(visitors, id_str)
		fmt.Println(visitors)

		book += 1
		fmt.Println("write book: " + strconv.Itoa(book))
		release(wsem)

		r = rand.Intn(max_time_to_sleep)
		time.Sleep(time.Duration(r) * time.Millisecond)

		visitors = remove_element(visitors, id_str)
		fmt.Println(visitors)
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