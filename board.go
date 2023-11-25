package main

import (
	"fmt"
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

func spawn_random_traveller(board [][]*vertex, traveller_type int) {
	switch traveller_type {
	case normal:
		for{
			traveller := new_traveller(k, traveller_type)
			k++
			go start_traveller(traveller, board)
			if k >= m * n { break }
			time.Sleep(normal_traveller_wait_time)
		}
		break
	case wild:
		for {
			// traveller := new_traveller(99, traveller_type)
			time.Sleep(wild_traveller_wait_time)
		}
		break
	case danger:
		for {
			// traveller := new_traveller(-1, traveller_type)
			time.Sleep(danger_traveller_wait_time)
		}
		break
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

func run_wild_locator(traveller *traveller, board [][]*vertex) {

}

func print_board(board [][]*vertex) {
	var i int
	var j int

	for {
		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				if board[i][j].traveller != nil{
					switch board[i][j].traveller.traveller_type {
					case normal:
						fmt.Printf("|%d", board[i][j].traveller.id)
						break
					case wild:
						fmt.Printf("|**")
						break
					case danger:
						fmt.Printf("|##")
						break
					}
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