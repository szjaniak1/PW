package main

import (
	"fmt"
	"time"
	"math/rand"
	"runtime"
)

var k int
var m int
var n int
var traces [][]uint8
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

			if vert.traveller.traveller_type == normal {
					read.resp <- false
					break
			}

			if vert.traveller.traveller_type == danger {
				read_op := read_op{
					action : you_die,
					resp : make(chan bool)}

				read.resp <- true
				new_traveller := <- vert.write_channel
				vert.traveller = nil
				new_traveller.resp <- true
				
				new_traveller.val.notify <- read_op
				break
			}

			switch read.action {
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
			default:
				read.resp <-false
				break
			}
		}
	}
}

func start_normal_traveller(traveller *traveller, board[][]*vertex, ws <-chan int) {
	read_op := read_op{
		action : normal_traveller_move_in,
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		pos_x := rand.Intn(m)
		pos_y := rand.Intn(n)

		board[pos_x][pos_y].read_channel <-read_op
		access := <-read_op.resp
		if !access { continue }

		board[pos_x][pos_y].write_channel <- write_op
		resp := <-write_op.resp
		if resp {
			traveller.pos_x = pos_x
			traveller.pos_y = pos_y
			break
		}
	}

	go run_normal_traveller(traveller, board, ws)
}

func start_wild_traveller(traveller *traveller, board[][]*vertex, ws <-chan int) {
	read_op := read_op{
		action : wild_traveller_move_in,
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		pos_x := rand.Intn(m)
		pos_y := rand.Intn(n)

		board[pos_x][pos_y].read_channel <-read_op
		access := <-read_op.resp
		if !access { continue }

		board[pos_x][pos_y].write_channel <- write_op
		resp := <-write_op.resp
		if resp {
			traveller.pos_x = pos_x
			traveller.pos_y = pos_y
			break
		}
	}

	go run_wild_traveller(traveller, board, ws)
}

func start_danger_traveller(traveller *traveller, board[][]*vertex, ws <-chan int) {
	read_op := read_op{
		action : danger_traveller_move_in,
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		pos_x := rand.Intn(m)
		pos_y := rand.Intn(n)

		board[pos_x][pos_y].read_channel <-read_op
		access := <-read_op.resp
		if !access { continue }

		board[pos_x][pos_y].write_channel <- write_op
		resp := <-write_op.resp
		if resp {
			traveller.pos_x = pos_x
			traveller.pos_y = pos_y
			break
		}
	}

	go run_wild_traveller(traveller, board, ws)
}

func run_danger_traveller(traveller *traveller, board [][]*vertex, ws <-chan int) {
	state := RUNNING
	timer := time.NewTimer(wild_traveller_life_time)
	quit := make(chan bool)


	go func() {
		for {
			select{
			case state = <-ws:
				switch state {
				case STOPPED:
					return
				case RUNNING:
				case PAUSED:
				}
			default:
				runtime.Gosched()

				if state == PAUSED {
					break
				}

				select{
				case <-quit:
					board[traveller.pos_x][traveller.pos_y].traveller = nil
					return
				}
			}
		}
	}()
	<-timer.C
	quit <- true
}

func run_normal_traveller(traveller *traveller, board [][]*vertex, ws <-chan int) {
	state := RUNNING

	read_op := read_op{
		action : normal_traveller_move_in,
		resp : make(chan bool)}
	write_op := write_op{
		val : traveller,
		resp : make(chan bool)}

	for {
		select{
		case state = <-ws:
			switch state {
			case STOPPED:
				return
			case RUNNING:
			case PAUSED:
			}
		case <- traveller.notify:
			traveller = nil
			return
		default:
			runtime.Gosched()

			if state == PAUSED {
				break
			}

			x := traveller.pos_x
			y := traveller.pos_y
			switch direction := rand.Intn(4); direction {
			case RIGHT:
				if x + 1 >= m { break }

				board[x + 1][y].read_channel <-read_op
				access := <-read_op.resp
				if !access { break }

				board[x + 1][y].write_channel <- write_op
				resp := <-write_op.resp
				if resp {
					traces[x][y] = 1
					traveller.pos_x++
					board[x][y].traveller = nil
				}
			break
			case LEFT:
				if x - 1 < 0 { break }

				board[x - 1][y].read_channel <- read_op
				access := <- read_op.resp
				if !access { break }

				board[x - 1][y].write_channel <- write_op
				resp := <- write_op.resp
				if resp {
					traces[x][y] = 1
					traveller.pos_x--
					board[x][y].traveller = nil
				}
				break
			case DOWN:
				if y + 1 >= n { break }

				board[x][y + 1].read_channel <- read_op
				access := <- read_op.resp
				if !access { break }

				board[x][y + 1].write_channel <- write_op
				resp := <- write_op.resp
				if resp {
					traces[x][y] = 1
					traveller.pos_y++
					board[x][y].traveller = nil
				}
				break
			case UP:
				if y - 1 < 0 { break }

				board[x][y - 1].read_channel <- read_op
				access := <- read_op.resp
				if !access { break }

				board[x][y - 1].write_channel <- write_op
				resp := <- write_op.resp
				if resp {
					traces[x][y] = 1
					traveller.pos_y--
					board[x][y].traveller = nil
				}
				break
			}
			duration := rand.Intn(normal_traveller_thinking_range)
			time.Sleep(time.Duration(duration) * time.Millisecond)
		}
	}
}

func run_wild_traveller(traveller *traveller, board [][]*vertex, ws <-chan int) {
	state := RUNNING
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
			select {
			case state = <-ws:
				switch state {
				case STOPPED:
					return
				case RUNNING:
				case PAUSED:
				}
			default:
				runtime.Gosched()

				if state == PAUSED {
					break
				}

				select {
				case read := <-traveller.notify:
					switch read.action {
					case you_die:
						traveller = nil
						return
					default:
						var access bool
						x := traveller.pos_x
						y := traveller.pos_y
						if x + 1 < m {
							board[x + 1][y].read_channel <- read_op
							access = <-read_op.resp
							if access {
								board[x + 1][y].write_channel <- write_op
								resp := <- write_op.resp
								if resp {
									traveller.pos_x++
									read.resp <- true
									break
								}
							}
						}
						
						if x - 1 >= 0 {
							board[x - 1][y].read_channel <- read_op
							access = <-read_op.resp
							if access {
								board[x - 1][y].write_channel <- write_op
								resp := <- write_op.resp
								if resp {
									traveller.pos_x--
									read.resp <- true
									break
								}
							}
						}
						
						if y + 1 < n {
							board[x][y + 1].read_channel <- read_op
							access = <-read_op.resp
							if access {
								board[x][y + 1].write_channel <- write_op
								resp := <- write_op.resp
								if resp {
									traveller.pos_y++
									read.resp <- true
									break
								}
							}
						}
						
						if y - 1 >= 0 {
							board[x][y - 1].read_channel <- read_op
							access = <-read_op.resp
							if access {
								board[x][y - 1].write_channel <- write_op
								resp := <- write_op.resp
								if resp {
									traveller.pos_y--
									read.resp <- true
									break
								}
							}	
						}
						read.resp <- false
						break
					}
				case <-quit:
					board[traveller.pos_x][traveller.pos_y].traveller = nil
					return
				}
			}
		}
	}()
	<-timer.C
	quit <- true
}

func create_travellers() []*traveller {
	var travellers []*traveller

	for i := 0; i < normal_limit; i++ {
		traveller := new_traveller(k, normal)
		travellers = append(travellers, traveller)
		k++
	}

	for i := 0; i < wild_limit; i++ {
		traveller := new_traveller(99, wild)
		travellers = append(travellers, traveller)
	}

	for i := 0; i < danger_limit; i++ {
		traveller := new_traveller(999, danger)
		travellers = append(travellers, traveller)
	}

	return travellers
}

func print_board(freezers []chan int, board [][]*vertex) {
	var i int
	var j int

	// height := 2 * m - 1
	// width := 2 * n - 1

	// printing_board := make([][]string, height)
	

	for {
		// for row_idx := in range printing_board {
		// 	row := make([]int, width)
		// 	printing_board[row_idx] = row
		// }

		// set_state(freezers, RUNNING)
		// for i = 0; i < m; i++ {
		// 	for j = 0; j < n; j++ {
		// 		if board[i][j].traveller == nil {
		// 			printing_board[i][j]
		// 			continue
		// 		}

		// 	}
		// }






		for i = 0; i < m; i++ {
			for j = 0; j < n; j++ {
				if board[i][j].traveller == nil{
					fmt.Printf("    ")
					continue
				}

				switch board[i][j].traveller.traveller_type {
				case normal:
					fmt.Printf("%d  ", board[i][j].traveller.id)
					break
				case wild:
					fmt.Printf("**  ")
					break
				case danger:
					fmt.Printf("##  ")
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
			fmt.Println()

			fmt.Println()
		}
		fmt.Println()
		set_state(freezers, RUNNING)
		time.Sleep(time.Second * 1)
	}
}

func set_state(travellers []chan int, state int) {
	for _, w := range travellers {
		go func() {
			w <- state
		}()
	}
}