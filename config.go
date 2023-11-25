package main

import "time"

const (
	normal = 0
	wild = 1
	danger = 2
)

const normal_traveller_wait_time = time.Second * 5
const wild_traveller_wait_time = time.Second * 7
const danger_traveller_wait_time = time.Second * 10