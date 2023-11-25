package main

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