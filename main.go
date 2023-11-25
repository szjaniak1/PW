package main

import(
	"fmt"
	"os"
	"strconv"
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

	// traces = make([][]uint8, m)
	// for i := range traces {
	// 	traces[i] = make([]uint8, n)
	// }

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

	go spawn_random_traveller(board, normal)
	go spawn_random_traveller(board, wild)
	go print_board(board)

	select {}
}