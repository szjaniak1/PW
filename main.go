package main

import(
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	argsWithoutProg := os.Args[1:]
	k = 10
	m, err = strconv.Atoi(argsWithoutProg[0])
	n, err = strconv.Atoi(argsWithoutProg[1])
	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	var board [][]*vertex

	for i := 0; i < m; i++ {
		var vertex_row []*vertex
		for j:= 0; j < n; j++ {
			ver := new_vertex()
			vertex_row = append(vertex_row, ver)
			go vertex_listener(ver)
		}
		board = append(board, vertex_row)
	}

	height = 2 * m - 1
	width = 2 * n - 1

	printing_board = make([][]string, height)
	for row := range printing_board {
		printing_board[row] = make([]string, width)
	}

	for i := 1; i < height; i += 2 {
		for j := 1; j < width; j += 2 {
			printing_board[i][j] = "  "
		}
	}


	var wg sync.WaitGroup
	worker_count := normal_limit + wild_limit + danger_limit
	wg.Add(worker_count + 1)
	freezers := make([]chan int, worker_count)
	travellers := create_travellers()

	go func() {
		for i := 0; i < normal_limit; i++ {
			go func(i int) {
				start_normal_traveller(travellers[i], board, freezers[i])
				wg.Done()
			}(i)
			time.Sleep(normal_traveller_wait_time)
		}
	}()

	go func() {
		for i := normal_limit; i < worker_count - danger_limit; i++ { //it should respawn them when they die after some time
			go func(i int) {
				start_wild_traveller(travellers[i], board, freezers[i])
				wg.Done()
			}(i)
			time.Sleep(wild_traveller_wait_time)
		}
	}()

	go func() {
		for i := worker_count - danger_limit; i < worker_count; i++ {
			go func(i int) {
				start_danger_traveller(travellers[i], board, freezers[i])
				wg.Done()
			}(i)
			time.Sleep(danger_traveller_wait_time)
		}
	}()
	
	time.Sleep(wild_traveller_wait_time)
	go func() {
		print_board(freezers, board)
		wg.Done()
	}()

	select {}
}