package main

import (
	"Mutual_Exclusion/m/v2/raalgo"
	"log"
	"net"
    "context"
    "sync"
    "fmt"

	"google.golang.org/grpc"
)

//Node needs to keep list of other nodes in network
//Look at map from ChittyChat - map of nodeid to stream


type P2PNode struct {
    pb.UnimplementedP2PServiceServer
    peers    map[string]pb.P2PServiceClient // map of peer addresses to clients
    peerLock sync.RWMutex
}

// Server method implementation
func (n *P2PNode) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
    log.Printf("Received message from %s to %s: %s\n", req.From, req.To, req.Message)
    return &pb.MessageResponse{Status: "Message received"}, nil
}

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

// Add a peer to the node
func (n *P2PNode) AddPeer(address string) {
    conn, err := grpc.Dial(address, grpc.WithInsecure())
    if err != nil {
        log.Printf("Failed to connect to peer %s: %v", address, err)
        return
    }

    client := pb.NewP2PServiceClient(conn)

    n.peerLock.Lock()
    n.peers[address] = client
    n.peerLock.Unlock()

    log.Printf("Connected to peer %s", address)
}

func main() {
    node := &P2PNode{
        peers: make(map[string]pb.P2PServiceClient),
    }

    // Start the gRPC server
    go func() {
        lis, err := net.Listen("tcp", ":50051")
        if err != nil {
            log.Fatalf("Failed to listen: %v", err)
        }

        grpcServer := grpc.NewServer()
        pb.RegisterP2PServiceServer(grpcServer, node)

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
        res, err := peer.SendMessage(context.Background(), &pb.MessageRequest{
            From:    "localhost:50051",
            To:      "localhost:50052",
            Message: "Hello, peer!",
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
