package traveller

import(
	"fmt"
	"time"
)

type traveller struct {
	id int
	pos_x int
	pos_y int
}

func new_traveller(id int) *traveller{
	traveller := traveller{ id : id }

	//here we need to place our traveller into the board or do it in run_traveller() method

	return &traveller
}

func run_traveller(traveller *traveller, board [][]*vertex) {
	read_op := make(chan read_op)
	write_op := make(chan write_op)
	write_op.val = traveller
	for {
		x := traveller.pos_x
		y := traveller.pos_y
		switch direction := rand.Intn(4); direction {
		case 0:
			if x + 1 < m {
				board[x + 1][y].read_channel <- read_op
				if access := <- read_op {
					board[x + 1][y].write_channel <- write_op
					if resp := <- write_op {
						//traces[x][y] = 1
						traveller.pos++
						board[x][y].traveller = nil
					}
				}
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