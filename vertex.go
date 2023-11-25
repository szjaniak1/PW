package main

type vertex struct {
	traveller *traveller
	read_channel chan read_op
	write_channel chan write_op
}

type read_op struct{
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