package main

type vertex struct {
	traveller *traveller
	read_channel chan read_op
	write_channel chan write_op
}

const (
	wild_traveller_move_in = 0
	wild_traveller_quit = 1
	normal_traveller_move_in = 2
)

type read_op struct{
	action int
	resp chan bool
}

type write_op struct{
	val *traveller
	resp chan bool
}

func new_vertex() *vertex{
	ver := vertex{ read_channel : make(chan read_op), write_channel : make(chan write_op) }

	return &ver
}