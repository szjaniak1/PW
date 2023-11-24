package main

import(
	"fmt"
	"os"
	"strconv"
	"time"
	"math/rand"
	"traveller"
	"vertex"
)

var k int
var m int
var n int
// var traces [][]uint8
var err error

func vertex_listener(vert *vertex) {
	for {
		select {
		case read := <- vert.read_channel:
			if vert.traveller {
				read.resp <- false
			}
			else {
				read.resp <- true
				new_traveller := <- vert.write_channel
				vert.traveller = new_traveller.val
				new_traveller.resp <- true 
			}
		}
	}
}

func spawn_random_traveller(board [][]*vertex) {
	for{
		// traveller := new_traveller(k, board)
		// k++
		// go run_traveller(traveller, board)
		// if k >= m * n { break }
		// time.Sleep(time.Second * 5)
	}
}

func print_board(board [][]*vertex) {
	var i int
	var j int

	// for {
	// 	for i = 0; i < m; i++ {
	// 		for j = 0; j < n; j++ {
	// 			if board[i][j].locator != "--" {
	// 				fmt.Printf("|%s", board[i][j].locator)
	// 				if traces[i][j] == 1 {
	// 					traces[i][j] = 0
	// 				}
	// 			} else if traces[i][j] == 1{
	// 				fmt.Printf("|xx")
	// 				traces[i][j] = 0
	// 				// printing traves should work differently
	// 			} else {
	// 				fmt.Printf("|--")
	// 			}
	// 		}
	// 		fmt.Printf("|")
	// 		fmt.Println()
	// 	}
	// 	fmt.Println()
	// 	time.Sleep(time.Second * 1)
	// }
}

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

	// var board [][]*vertex //our board should be global for the single vertexes to be able to look at the others

	// for i := 0; i < m; i++ {
	// 	var vertex_row []*vertex
	// 	for j:= 0; j < n; j++ {
	// 		ver := new_vertex("--", make(chan *traveller))
	// 		vertex_row = append(vertex_row, ver)
	// 		go vertex_listener(ver)
	// 	}
	// 	board = append(board, vertex_row)
	// }

	// go spawn_random_traveller(board)
	// go print_board(board)

	select {}
}