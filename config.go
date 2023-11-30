package main

import "time"

const normal_limit = 10
const wild_limit = 10
const danger_limit = 20

const (
	normal = 0
	wild = 1
	danger = 2
)

const normal_traveller_wait_time = time.Second * 2
const normal_traveller_thinking_range = 500
const wild_traveller_wait_time = time.Second * 3
const wild_traveller_life_time = time.Second * 7
const danger_traveller_wait_time = time.Second * 3
const danger_traveller_life_time = time.Second * 7

const (
	wild_traveller_move_in		= 0
	normal_traveller_move_in	= 1
	danger_traveller_move_in	= 2
	you_die						= 3
)

const (
	RIGHT 	= 0
	LEFT 	= 1
	DOWN	= 2
	UP 		= 3
)

const (
	STOPPED = 0
	PAUSED 	= 1
	RUNNING = 2
)