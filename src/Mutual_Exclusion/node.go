package main

import (
	pb "Mutual_Exclusion/m/v2/raalgo"
	"context"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//Node needs to keep list of other nodes in network
//Look at map from ChittyChat - map of nodeid to stream

type P2PNode struct {
	pb.UnimplementedP2PnetworkServer
	peers             map[string]pb.P2PnetworkClient // map of peer addresses to clients
	peerLock          sync.RWMutex
	ME                int64
	N                 int64
	Our_Timestamp     int64
	Highest_Timestamp int64
	Outstanding_Reply int64
	Request_Critical  bool
	Reply_Defered     []bool
}

// Server method implementation
func (n *P2PNode) SendMessage(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
	log.Printf("Received message from %d at time %d\n", req.Nodeid, req.Timestamp)
	return &pb.Reply{Permission: true}, nil
}

//Below is an implementation of a GetStatus rpc method that we don't have in our .proto
//Returns which nodes a node is currently connected to
//Keep it for eventual debugging purposes
/**
func (n *P2PNode) GetStatus(ctx context.Context, req *pb.Empty) (*pb.StatusResponse, error) {
    n.peerLock.RLock()
    defer n.peerLock.RUnlock()

    var peers []string
    for addr := range n.peers {
        peers = append(peers, addr)
    }

    return &pb.StatusResponse{
        Status:        "Active",
        ConnectedPeers: peers,
    }, nil
}
*/

// Add a peer to the node
func (n *P2PNode) AddPeer(address string) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to peer %s: %v", address, err)
		return
	}

	client := pb.NewP2PnetworkClient(conn)

	n.peerLock.Lock()
	n.peers[address] = client
	n.peerLock.Unlock()

	log.Printf("Connected to peer %s", address)
}

func StartServer(node *P2PNode, address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterP2PnetworkServer(grpcServer, node)

	log.Printf("P2P node is running on port %s/n", address)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func CreateNode(l int64, n int64) *P2PNode {
	node := &P2PNode{
		ME:                l,
		N:                 n,
		Request_Critical:  false,
		peers:             make(map[string]pb.P2PnetworkClient),
		Our_Timestamp:     1,
		Highest_Timestamp: 1,
		Outstanding_Reply: n - 1,
		Reply_Defered:     make([]bool, n),
	}
	for i := 0; i < len(node.Reply_Defered); i++ {
		node.Reply_Defered[i] = false
	}
	return node
}

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
	node1.AddPeer("localhost:50052") // Example of connecting to another node
	node1.AddPeer("localhost:50053")

	node2.AddPeer("localhost:50051") // Example of connecting to another node
	node2.AddPeer("localhost:50053")

	node3.AddPeer("localhost:50051") // Example of connecting to another node
	node3.AddPeer("localhost:50052")

	// Simulate sending a message to a peer
	node1.peerLock.RLock()
	if peer, exists := node1.peers["localhost:50052"]; exists {
		res, err := peer.Ask(context.Background(), &pb.Request{
			Nodeid:    1, //pass in
			Timestamp: 2, //keep track
		})
		if err != nil {
			log.Printf("Error sending message: %v", err)
		} else {
			if res.Permission {
				log.Print("reply received")
			}
		}
	}
	node1.peerLock.RUnlock()

	for {
		node1.Request_Critical = true
		log.Printf("Node 1 requests access \n")
		node1.Ask(context.Background(), &pb.Request{
			Nodeid:    int64(node1.ME),
			Timestamp: int64(node1.Our_Timestamp + 1),
		})

		while
	}
}
