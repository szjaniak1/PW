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
			go func() {
				if vert.traveller != nil{
					if vert.traveller.traveller_type == wild {
						switch read.action {
						case wild_traveller_move_in:
							read.resp <- false
							break
						case normal_traveller_move_in:
							// fmt.Printf("JESTESMYYY %d;%d\n", vert.traveller.pos_x, vert.traveller.pos_y)
							read_op := read_op{
								action : normal_traveller_move_in,
								resp : make(chan bool)}
							vert.traveller.notify <- read_op
							// fmt.Printf("JESTESMYYY moze to %d;%d\n", vert.traveller.pos_x, vert.traveller.pos_y)
							access := <- read_op.resp
							// fmt.Printf("JESTESMYYY???? %d;%d\n", vert.traveller.pos_x, vert.traveller.pos_y)
							if access {
								read.resp <- true
								new_traveller := <- vert.write_channel
								vert.traveller = new_traveller.val
								new_traveller.resp <- true 
							} else {
								read.resp <- false
							}
							break
						case wild_traveller_quit:
							read.resp <- true
							new_traveller := <- vert.write_channel
							vert.traveller = new_traveller.val
							new_traveller.resp <- true 
							break
						}
					} else {
						read.resp <- false
					}
				} else {
					// fmt.Printf("empty 4\n")
					read.resp <- true
					new_traveller := <- vert.write_channel
					vert.traveller = new_traveller.val
					// fmt.Printf("empty %d;%d\n", vert.traveller.pos_x, vert.traveller.pos_y)
					new_traveller.resp <- true 
				}
			}()
		}
	}
}

func spawn_random_traveller(board [][]*vertex, traveller_type int) {
	limit := m * n
	switch traveller_type {
	case normal:
		for{
			traveller := new_traveller(k, traveller_type)
			k++
			go start_traveller(traveller, board, traveller_type)
			if k >= limit { break }
			time.Sleep(normal_traveller_wait_time)
		}
		break
	case wild:
		for {
			traveller := new_traveller(99, traveller_type)
			go start_traveller(traveller, board, traveller_type)
			if k >= limit{ break }
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
	switch traveller_type {
	case normal:
		go run_traveller(traveller, board)
		break
	case wild:
		go run_wild_traveller(traveller, board)
		break
	}
}

func run_traveller(traveller *traveller, board [][]*vertex) {
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
				// fmt.Printf("start %d;%d\n", x, y)
				if x + 1 < m {
					// fmt.Printf("check 11 %d;%d\n", x, y)
					board[x + 1][y].read_channel <- read_op
					// fmt.Printf("check 1 %d;%d\n", x, y)
					access = <-read_op.resp
					// fmt.Printf("check 444 %d;%d\n", x, y)
					if access {
						board[x + 1][y].write_channel <- write_op
						resp := <- write_op.resp
						if resp {
							read.resp <- true
						}
						break
					}
				}
				
				if x - 1 >= 0 {
					// fmt.Printf("check 22 %d;%d\n", x, y)
					board[x - 1][y].read_channel <- read_op
					// fmt.Printf("check 2 %d;%d\n", x, y)
					access = <-read_op.resp
					// fmt.Printf("check 222 %d;%d\n", x, y)
					if access {
						board[x - 1][y].write_channel <- write_op
						resp := <- write_op.resp
						if resp {
							read.resp <- true
						}
						break
					}
				}
				
				if y + 1 < n {
					// fmt.Printf("check 33 %d;%d\n", x, y)
					board[x][y + 1].read_channel <- read_op
					// fmt.Printf("check 3 %d;%d\n", x, y)
					access = <-read_op.resp
					// fmt.Printf("check 333 %d;%d\n", x, y)
					if access {
						board[x][y + 1].write_channel <- write_op
						
						resp := <- write_op.resp
						if resp {
							read.resp <- true
						}
						break
					}
				}
				
				if y - 1 >= 0 {
					// fmt.Printf("check 44 %d;%d\n", x, y)
					board[x][y - 1].read_channel <- read_op
					// fmt.Printf("check 4 %d;%d\n", x, y)
					access = <-read_op.resp
					// fmt.Printf("check 444 %d;%d\n", x, y)
					if access {
						board[x][y - 1].write_channel <- write_op

						resp := <- write_op.resp
						if resp {
							read.resp <- true
						}
						break
					}
				}
				// fmt.Printf("cannot %d;%d\n", x, y)
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