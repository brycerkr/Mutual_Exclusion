package main

import (
	"time"
)

func main() {
	node1 := CreateNode(0, 3)
	node2 := CreateNode(1, 3)
	node3 := CreateNode(2, 3)

	// Start the gRPC server
	go StartServer(node1, ":50051")
	go StartServer(node2, ":50052")
	go StartServer(node3, ":50053")

	time.Sleep(3 * time.Second)

	// Add peers (simulate peer discovery for demonstration)
	node1.AddPeer("localhost:50052")
	node1.AddPeer("localhost:50053")

	node2.AddPeer("localhost:50051")
	node2.AddPeer("localhost:50053")

	node3.AddPeer("localhost:50051")
	node3.AddPeer("localhost:50052")

	go node1.Start()
	go node2.Start()
	go node3.Start()
}
