package main

import "time"

const (
	normal = 0
	wild = 1
	danger = 2
)

const normal_traveller_wait_time = time.Second * 3
const wild_traveller_wait_time = time.Second * 4
const wild_traveller_life_time = time.Second * 6
const danger_traveller_wait_time = time.Second * 10

const (
	wild_traveller_move_in		= 0
	wild_traveller_quit			= 1
	normal_traveller_move_in	= 2
)

const (
	RIGHT 	= 0
	LEFT 	= 1
	DOWN	= 2
	UP 		= 3
)