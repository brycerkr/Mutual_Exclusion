package main

var balance int16

func resource() {
	balance++
}

func CriticalSection() {
	resource()
}
