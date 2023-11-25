package main

type traveller struct {
	id int
	traveller_type int
	pos_x int
	pos_y int
}

func new_traveller(id int, traveller_type int) *traveller{
	traveller := traveller{ id : id, traveller_type : traveller_type }

	return &traveller
}