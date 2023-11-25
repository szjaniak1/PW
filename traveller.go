package main

type traveller struct {
	id int
	pos_x int
	pos_y int
}

func new_traveller(id int) *traveller{
	traveller := traveller{ id : id }

	return &traveller
}