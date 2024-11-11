package main

import (
	"log"
	"time"
)

var balance = 0

func resource() {
	balance++
	log.Printf("Balance: %d", balance)
}

func CriticalSection() {
	resource()
}
