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
// var traces [][]uint8
var err error

func vertex_listener(vert *vertex) {
	for {
		select {
		case read := <- vert.read_channel:
			if vert.traveller != nil{
				read.resp <- false
			} else {
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
		traveller := new_traveller(k)
		k++
		go start_traveller(traveller, board)
		if k >= m * n { break }
		time.Sleep(time.Second * 5)
	}
}

func start_traveller(traveller *traveller, board[][]*vertex) {
	read_op := read_op{
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		pos_x := rand.Intn(m)
		pos_y := rand.Intn(n)

		board[pos_x][pos_y].read_channel <-read_op
		access := <-read_op.resp
		if access {
			board[pos_x][pos_y].write_channel <- write_op
			resp := <-write_op.resp
			if resp {
				traveller.pos_x = pos_x
				traveller.pos_y = pos_y
				break
			}
		}
	}
	go run_traveller(traveller, board)
}

func run_traveller(traveller *traveller, board [][]*vertex) {
	read_op := read_op{
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		x := traveller.pos_x
		y := traveller.pos_y
		switch direction := rand.Intn(4); direction {
		case 0:
			if x + 1 < m {
				board[x + 1][y].read_channel <-read_op
				access := <-read_op.resp
				if access {
					board[x + 1][y].write_channel <- write_op
					resp := <-write_op.resp
					if resp {
						//traces[x][y] = 1
						traveller.pos_x++
						board[x][y].traveller = nil
					}
				}
			}
			break
		case 1:
			if x - 1 >= 0{
				board[x - 1][y].read_channel <- read_op
				access := <- read_op.resp
				if access {
					board[x - 1][y].write_channel <- write_op
					resp := <- write_op.resp
					if resp {
						//traces[x][y] = 1
						traveller.pos_x--
						board[x][y].traveller = nil
					}
				}
			}
			break
		case 2:
			if y + 1 < n {
				board[x][y + 1].read_channel <- read_op
				access := <- read_op.resp
				if access {
					board[x][y + 1].write_channel <- write_op
					resp := <- write_op.resp
					if resp {
						//traces[x][y] = 1
						traveller.pos_y++
						board[x][y].traveller = nil
					}
				}
			}
			break
		case 3:
			if y - 1 >= 0 {
				board[x][y - 1].read_channel <- read_op
				access := <- read_op.resp
				if access {
					board[x][y - 1].write_channel <- write_op
					resp := <- write_op.resp
					if resp {
						//traces[x][y] = 1
						traveller.pos_y--
						board[x][y].traveller = nil
					}
				}
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
				if board[i][j].traveller != nil{
					fmt.Printf("|%d", board[i][j].traveller.id)
					// if traces[i][j] == 1 {
					// 	traces[i][j] = 0
					// }
				// } else if traces[i][j] == 1{
				// 	fmt.Printf("|xx")
				// 	traces[i][j] = 0
				// 	// printing traves should work differently
				} else {
					fmt.Printf("|  ")
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

	go spawn_random_traveller(board)
	go print_board(board)

	select {}
}