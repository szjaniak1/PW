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
			if vert.traveller == nil{
				read.resp <- true
				new_traveller := <- vert.write_channel
				vert.traveller = new_traveller.val
				new_traveller.resp <- true
				break
			}

			if vert.traveller.traveller_type != wild {
					read.resp <- false
					break
			}

			switch read.action {
				case wild_traveller_move_in:
					read.resp <- false
					break
				case normal_traveller_move_in:
					read_op := read_op{
						action : normal_traveller_move_in,
						resp : make(chan bool)}
					vert.traveller.notify <- read_op
					access := <- read_op.resp
					if !access {
						read.resp <- false
						break
					}
					
					read.resp <- true
					new_traveller := <- vert.write_channel
					vert.traveller = new_traveller.val
					new_traveller.resp <- true 
					break
				case wild_traveller_quit:
					read.resp <- true
					new_traveller := <- vert.write_channel
					vert.traveller = new_traveller.val
					new_traveller.resp <- true 
					break
			}
		}
	}
}

func spawn_normal_traveller(board [][]*vertex) {
	limit := m * n
	for{
		traveller := new_traveller(k, normal)
		k++
		go start_traveller(traveller, board, normal)
		if k >= limit { break }
		time.Sleep(normal_traveller_wait_time)
	}
}

func spawn_wild_traveller(board [][]*vertex) {
	limit := m * n
	for{
		traveller := new_traveller(k, wild)
		k++
		go start_traveller(traveller, board, wild)
		if k >= limit { break }
		time.Sleep(wild_traveller_wait_time)
	}
}

func spawn_random_traveller(board [][]*vertex, traveller_type int) {
	switch traveller_type {
	case normal:
		spawn_normal_traveller(board)
		break
	case wild:
		spawn_wild_traveller(board)
		break
	case danger:
		for {
			time.Sleep(danger_traveller_wait_time)
		}
		break
	}
}

func start_traveller(traveller *traveller, board[][]*vertex, traveller_type int) {
	var action int
	if traveller_type == normal {
		action = normal_traveller_move_in
	} else if traveller_type == wild {
		action = wild_traveller_move_in
	}
	read_op := read_op{
		action : action,
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		pos_x := rand.Intn(m)
		pos_y := rand.Intn(n)

		board[pos_x][pos_y].read_channel <-read_op
		access := <-read_op.resp
		if !access { break }

		board[pos_x][pos_y].write_channel <- write_op
		resp := <-write_op.resp
		if resp {
			traveller.pos_x = pos_x
			traveller.pos_y = pos_y
			break
		}
	}
	switch traveller_type {
	case normal:
		go run_normal_traveller(traveller, board)
		break
	case wild:
		go run_wild_traveller(traveller, board)
		break
	}
}

func run_normal_traveller(traveller *traveller, board [][]*vertex) {
	read_op := read_op{
		action : normal_traveller_move_in,
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		x := traveller.pos_x
		y := traveller.pos_y
		switch direction := rand.Intn(4); direction {
		case 0:
			if x + 1 >= m { break }

			board[x + 1][y].read_channel <-read_op
			access := <-read_op.resp
			if !access { break }

			board[x + 1][y].write_channel <- write_op
			resp := <-write_op.resp
			if resp {
				//traces[x][y] = 1
				traveller.pos_x++
				board[x][y].traveller = nil
			}
		break
		case 1:
			if x - 1 < 0 { break }

			board[x - 1][y].read_channel <- read_op
			access := <- read_op.resp
			if !access { break }

			board[x - 1][y].write_channel <- write_op
			resp := <- write_op.resp
			if resp {
				//traces[x][y] = 1
				traveller.pos_x--
				board[x][y].traveller = nil
			}
			break
		case 2:
			if y + 1 >= n { break }

			board[x][y + 1].read_channel <- read_op
			access := <- read_op.resp
			if !access { break }

			board[x][y + 1].write_channel <- write_op
			resp := <- write_op.resp
			if resp {
				//traces[x][y] = 1
				traveller.pos_y++
				board[x][y].traveller = nil
			}
			break
		case 3:
			if y - 1 < 0 { break }

			board[x][y - 1].read_channel <- read_op
			access := <- read_op.resp
			if !access { break }

			board[x][y - 1].write_channel <- write_op
			resp := <- write_op.resp
			if resp {
				//traces[x][y] = 1
				traveller.pos_y--
				board[x][y].traveller = nil
			}
			break
		}
		duration := rand.Intn(500)
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}
}

func run_wild_traveller(traveller *traveller, board [][]*vertex) {
	timer := time.NewTimer(wild_traveller_life_time)
	quit := make(chan bool)

	read_op := read_op{
		action : wild_traveller_move_in,
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}
	go func() {
		for {
			select{
			case <-quit:
				go delete_wild_traveller(traveller, board)
				return
			case read := <-traveller.notify:
				var access bool
				x := traveller.pos_x
				y := traveller.pos_y
				if x + 1 < m {
					board[x + 1][y].read_channel <- read_op
					access = <-read_op.resp
					if !access { break }

					board[x + 1][y].write_channel <- write_op
					resp := <- write_op.resp
					if resp {
						read.resp <- true
					}
					break
				}
				
				if x - 1 >= 0 {
					board[x - 1][y].read_channel <- read_op
					access = <-read_op.resp
					if !access { break }

					board[x - 1][y].write_channel <- write_op
					resp := <- write_op.resp
					if resp {
						read.resp <- true
					}
					break
				}
				
				if y + 1 < n {
					board[x][y + 1].read_channel <- read_op
					access = <-read_op.resp
					if access { break }

					board[x][y + 1].write_channel <- write_op
					resp := <- write_op.resp
					if resp {
						read.resp <- true
					}
					break
				}
				
				if y - 1 >= 0 {
					board[x][y - 1].read_channel <- read_op
					access = <-read_op.resp
					if access { break }

					board[x][y - 1].write_channel <- write_op
					resp := <- write_op.resp
					if resp {
						read.resp <- true
					}
					break
				}
				read.resp <- false
			}
		}
	}()
	<-timer.C
	quit <- true
}

func delete_wild_traveller(traveller *traveller, board[][]*vertex) {
	x := traveller.pos_x
	y := traveller.pos_y

	read_op := read_op{
		action : wild_traveller_quit,
		resp : make(chan bool)}
	write_op := write_op{
		val : nil,
		resp : make(chan bool)}

	board[x][y].read_channel <- read_op
	access := <-read_op.resp
	if access {
		board[x][y].write_channel <- write_op
		<- write_op.resp
	}

	read := <-traveller.notify
	read.resp <- true
	close(traveller.notify)
}

func print_board(board [][]*vertex) {
	var i int
	var j int

	for {
		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				if board[i][j].traveller == nil{
					fmt.Printf("|  ")
					continue
				}

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
			}
			fmt.Printf("|")
			fmt.Println()
		}
		fmt.Println()
		time.Sleep(time.Second * 1)
	}
}