package main

import (
	"log"
	"time"
)

var balance = 0

func resource() {
	balance++
	log.Print(balance)
	time.Sleep(2 / time.Second)
}

func CriticalSection() {
	resource()
}
