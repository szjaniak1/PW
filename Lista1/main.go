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

func new_vertex(symbol string, channel chan *traveller) *vertex{
	ver := vertex{ locator : symbol, channel : channel }

	return &ver
}

func get_random_empty_vertex(trav *traveller, board [][]*vertex) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	m1 := r1.Intn(m)
	n1 := r1.Intn(n)

	if board[m1][n1].channel == nil {
		board[m1][n1].channel <- trav
		board[m1][n1].locator = strconv.Itoa(trav.num)
		trav.pos_x = m1
		trav.pos_y = n1
		return
	}

	get_random_empty_vertex(trav, board)
}

func vertex_listener(vert *vertex) {
	for {
		traveller := <- vert.channel
		vert.locator = strconv.Itoa(traveller.num)
	}
}

type traveller struct {
	num int
	pos_x int
	pos_y int
}

func new_traveller(num int, board [][]*vertex) *traveller{
	traveller := traveller{ num : num }

	get_random_empty_vertex(&traveller, board)

	return &traveller
}

func spawn_random_traveller(board [][]*vertex) {
	for{
		traveller := new_traveller(k, board)
		k++
		go run_traveller(traveller, board)
		if k >= m * n { break }
		time.Sleep(time.Second * 10)
	}
}

func run_traveller(traveller *traveller, board [][]*vertex) {
	for {
		x := traveller.pos_x
		y := traveller.pos_y
		switch direction := rand.Intn(4); direction {
		case 0:
			if x + 1 < m && board[x + 1][y].channel == nil{
				board[x + 1][y].channel <- traveller
				board[x][y].channel = nil
			}
		case 1:
			if x - 1 >= 0 && board[x - 1][y].channel == nil{
				board[x - 1][y].channel <- traveller
				board[x][y].channel = nil
			}
		case 2:
			if y + 1 < n && board[x][y + 1].channel == nil{
				board[x][y + 1].channel <- traveller
				board[x][y].channel = nil
			}
		case 3:
			if y >= 0 && board[x][y - 1].channel == nil{
				board[x][y - 1].channel <- traveller
				board[x][y].channel = nil
			}

		duration := rand.Intn(500)
		time.Sleep(time.Duration(duration) * time.Millisecond)
		}
	}
}

func print_board(board [][]*vertex) {
	var i int
	var j int

	for {
		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				if board[i][j].channel == nil {
					fmt.Printf("|--|")
				} else {
					fmt.Printf(board[i][j].locator)
				}
			}
			fmt.Println()
		}

		time.Sleep(time.Second * 3)
	}
}

func main() {
	argsWithoutProg := os.Args[1:]
	k = 1
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

	go spawn_random_traveller(board)
	go print_board(board)

	select {}
}