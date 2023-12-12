package main

type traveller struct {
	id int
	traveller_type int
	pos_x int
	pos_y int
	notify chan read_op
}

func new_traveller(id int, traveller_type int) *traveller{
	traveller := traveller{ id : id, traveller_type : traveller_type, notify : make(chan read_op) }

	return &traveller
}