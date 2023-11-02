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
var traces [][]uint8
var err error

type vertex struct {
	locator string
	traveller *traveller
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

	if board[m1][n1].traveller == nil {
		board[m1][n1].channel <- trav
		trav.pos_x = m1
		trav.pos_y = n1
		return
	}

	get_random_empty_vertex(trav, board)
}

func vertex_listener(vert *vertex) {
	for {
		traveller := <- vert.channel
		vert.traveller = traveller
		vert.locator = strconv.Itoa(traveller.id)

		if traveller.vertex != nil {
			traveller.vertex.locator = "--"
			traveller.vertex.traveller = nil
		}
		traveller.vertex = vert
	}
}

type traveller struct {
	id int
	pos_x int
	pos_y int
	vertex *vertex
}

func new_traveller(id int, board [][]*vertex) *traveller{
	traveller := traveller{ id : id }

	get_random_empty_vertex(&traveller, board)

	return &traveller
}

func spawn_random_traveller(board [][]*vertex) {
	for{
		traveller := new_traveller(k, board)
		k++
		go run_traveller(traveller, board)
		if k >= m * n { break }
		time.Sleep(time.Second * 5)
	}
}

func run_traveller(traveller *traveller, board [][]*vertex) {
	for {
		x := traveller.pos_x
		y := traveller.pos_y
		switch direction := rand.Intn(4); direction {
		case 0:
			if x + 1 < m && board[x + 1][y].traveller == nil{
				board[x + 1][y].channel <- traveller
				traveller.pos_x++
				traces[x][y] = 1
			}
			break
		case 1:
			if x - 1 >= 0 && board[x - 1][y].traveller == nil{
				board[x - 1][y].channel <- traveller
				traveller.pos_x--
				traces[x][y] = 1
			}
			break
		case 2:
			if y + 1 < n && board[x][y + 1].traveller == nil{
				board[x][y + 1].channel <- traveller
				traveller.pos_y++
				traces[x][y] = 1
			}
			break
		case 3:
			if y - 1 >= 0 && board[x][y - 1].traveller == nil{
				board[x][y - 1].channel <- traveller
				traveller.pos_y--
				traces[x][y] = 1
			}
			break
		}

		duration := rand.Intn(500)
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}
}

func print_board(board [][]*vertex) {
	var i int
	var j int

	for {
		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				if board[i][j].locator != "--" {
					fmt.Printf("|%s", board[i][j].locator)
					if traces[i][j] == 1 {
						traces[i][j] = 0
					}
				} else if traces[i][j] == 1{
					fmt.Printf("|xx")
					traces[i][j] = 0
				} else {
					fmt.Printf("|--")
				}
			}
			fmt.Printf("|")
			fmt.Println()
		}
		fmt.Println()
		time.Sleep(time.Second * 1)
	}
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

	traces = make([][]uint8, m)
	for i := range traces {
		traces[i] = make([]uint8, n)
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