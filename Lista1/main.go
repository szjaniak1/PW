package main

import(
	"fmt"
	"os"
	"strconv"
	"time"
	"math/rand"
)

var k int
var m int
var n int
var err error

type vertex struct {
	locator string
	channel chan *traveller
}

func new_vertex(symbol string, channel chan *traveller) {
	ver := vertex{
		locator : symbol,
		channel : channel
	}

	return &ver
}

func get_random_empty_vertex(traveller *traveller) {
	s1 := rand.NewSource(time.now().UnixNano())
	r1 := rand.New(s1)

	m1 := r1.Intn(m)
	n1 := r1.Intn(n)

	if (!board[m1][n1].traveller) 
		board[m1][n1].traveller <- traveller
		traveller.pos_x = m1
		traveller.pos_y = n1
		return

	get_random_empty_vertex(traveller)
}

func vertex_listener(vert *vertex) {
	for {
		traveller := <- vert.channel
		vert.traveller = traveller
		vert.locator = strconv.Itoa(traveller.num)
	}
}

type traveller struct {
	num int16
	pos_x int16
	pos_y int16
}

func new_traveller(num int) {
	traveller := traveller{
		num : num
	}

	get_random_empty_vertex(traveller)

	return &traveller
}

func spawn_random_traveller() {
	for{
		traveller := new_traveller(++k)
		go run_traveller(traveller)
		if k >= m * n break
		time.Sleep(time.Second * 10)
	}
}

func run_traveller(traveller *traveller, board [][]*vertex) {
	for {
		x := traveller.pos_x
		y := traveller.pos_y
		switch direction := rand.Intn(4); direction {
		case 0:
			if x + 1 < m && !board[x + 1][y].traveller
				board[x + 1][y].channel <= traveller
				board[x][y].traveller = nil
		case 1:
			if x - 1 >= 0 && !board[x - 1][y].traveller
				board[x - 1][y].channel <= traveller
				board[x][y].traveller = nil
		case 2:
			if y + 1 < n && !board[x][y + 1].traveller
				board[x][y + 1].channel <= traveller
				board[x][y].traveller = nil
		case 3:
			if y >= 0 && !board[x][y - 1].traveller
				board[x][y - 1].channel <= traveller
				board[x][y].traveller = nil

		duration := rand.Intn(500)
		time.Sleep(time.Duration(duration) * time.Milisecond)
		}
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

	var board [][]*vertex

	for i := 0; i < m; i++ {
		var vertex_row []*vertex
		for j:= 0; j < n; j++ {
			ver := new_vertex("--", make(chan *traveller))
			vertex_row = append(vertex_row, ver)
			go vertex_listener(ver)
		}
		board = append(board, vertex_row)
	}

	go spawn_random_traveller()
	go print_board(board)
	time.Sleep(time.Second* 5)
	board[2][2] = 45
	board[1][1] = 12
	time.Sleep(time.Second * 10)
}