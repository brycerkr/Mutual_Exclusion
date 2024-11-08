package main

import (
	pb "Mutual_Exclusion/m/v2/raalgo"
	"context"
	"fmt"
	"log"
	"net"
	"sync"

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
func (n *P2PNode) SendMessage(ctx context.Context, req *pb.Request) (*pb.Reply, error) {
    log.Printf("Received message from %d at time %d\n", req.Nodeid, req.Timestamp)
    return &pb.Reply{Permission: true}, nil
}

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

func main() {
    node := &P2PNode{
        peers: make(map[string]pb.P2PnetworkClient),
    }

    // Start the gRPC server
    go func() {
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
            log.Fatalf("Failed to listen: %v", err)
        }

        grpcServer := grpc.NewServer()
        pb.RegisterP2PnetworkServer(grpcServer, node)

        log.Println("P2P node is running on port 50051")
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("Failed to serve: %v", err)
        }
    }()

    // Add peers (simulate peer discovery for demonstration)
    node.AddPeer("localhost:50052") // Example of connecting to another node
    node.AddPeer("localhost:50053")

    // Simulate sending a message to a peer
    node.peerLock.RLock()
    if peer, exists := node.peers["localhost:50052"]; exists {
        res, err := peer.Ask(context.Background(), &pb.Request{
            Nodeid: 1,
            Timestamp: timestamp,
        })
        if err != nil {
            log.Printf("Error sending message: %v", err)
        } else {
            log.Printf("Response from peer: %s", res.Status)
        }
    }
    node.peerLock.RUnlock()

    select {} // Keep the main function running
}
