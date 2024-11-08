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
	peers    map[string]pb.P2PnetworkClient // map of peer addresses to clients
	peerLock sync.RWMutex
}

// Server method implementation
func (n *P2PNode) Ask(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
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

func main() {
	node1 := &P2PNode{
		peers: make(map[string]pb.P2PnetworkClient),
	}
    node2 := &P2PNode{
		peers: make(map[string]pb.P2PnetworkClient),
	}
    node3 := &P2PNode{
		peers: make(map[string]pb.P2PnetworkClient),
	}

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

	select {} // Keep the main function running
}
