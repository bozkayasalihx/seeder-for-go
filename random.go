package main

import "math/rand"

func randomId(len int) int {
	randomId := rand.Intn(len) + 1
	return randomId
}
