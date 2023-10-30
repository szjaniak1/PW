package main

import(
	"fmt"
	"os"
	"strconv"
	"time"
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
		time.Sleep(time.Second * 10)

		// traveler new
		go run_traveller(traveller)
	}	
}

func run_traveller(traveller *traveller) {
	for {

	}
}

func print_board(board [][]int) {
	var i int
	var j int

	for {
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

		time.Sleep(time.Second * 3)
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

	go spawn_random_traveller()
	go print_board(board)
	time.Sleep(time.Second* 5)
	board[2][2] = 45
	board[1][1] = 12
	time.Sleep(time.Second * 10)
}