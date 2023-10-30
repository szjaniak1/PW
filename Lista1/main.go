package main

import(
	"fmt"
	"os"
	"strconv"
)

var k int
var m int
var n int
var err error

type traveller struct {
	num int16
	pos_x int16
	pos_y int16
}

func spawn_random_traveller() {
	for{

	}
}

func run_traveller() {
}

func print_board(board [][]int) {
	var i int
	var j int
	for i = 0; i < m; i++ {
		for j = 0; j < n; j++ {
			if board[i][j] == 0 {
				fmt.Printf("|--|")
			} else {
				fmt.Printf("|%d|", board[i][j])
			}
		}
		fmt.Println()
	}
}

func generate_board(board [][]int) {
	var i int
	var j int
	// board = make([][]int, int(m), int(n))
	for i = 0; i < m; i++ {
		for j = 0; j < n; j++ {
			board[i][j] = 0
		}
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	k = 0
	m, err = strconv.Atoi(argsWithoutProg[0])
	n, err = strconv.Atoi(argsWithoutProg[1])

	if err != nil {
		fmt.Println("Error during conversion")
		return
	}

	var board [][]int = make([][]int, m)
	for i := range board {
		board[i] = make([]int, n)
	}

	generate_board(board)

	fmt.Printf(string(m + n + k))

	print_board(board)
	// go spawn_random_traveller()
}