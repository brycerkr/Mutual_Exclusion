package main

import "time"

var balance int16

func resource() {
	balance++
	time.Sleep(2 / time.Second)
}

func CriticalSection() {
	resource()
}
