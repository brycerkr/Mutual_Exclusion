package main

import (
	"Mutual_Exclusion/m/v2/raalgo"
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
func (n *P2PNode) Ask(ctx context.Context, req *pb.Request) (*raalgo.Reply, error) {
	log.Printf("Received message from %d at time %d\n", req.Nodeid, req.Timestamp)

	Defer_Request := false
	n.Highest_Timestamp = max(n.Highest_Timestamp, req.Timestamp)
	Defer_Request = n.Request_Critical && ((req.Timestamp > n.Our_Timestamp) ||
		(req.Timestamp == n.Our_Timestamp && req.Nodeid > n.ME))
	if Defer_Request {
		n.Reply_Defered[req.Nodeid] = true
	} else {
		//Answer(Permission)
	}

	return nil, nil
	// Reply needs to be deferred until Critical Section has been accessed
	// Need to implement CS logic

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

// Calls Ask with a Request for all peers that the node knows
func (n *P2PNode) AskAllPeers() {
	n.peerLock.RLock()
	for address := range n.peers {
		if peer, exists := n.peers[address]; exists {
			res, err := peer.Ask(context.Background(), &pb.Request{
				Nodeid:    n.ME,            //pass in
				Timestamp: n.Our_Timestamp, //keep track
			})
			if err != nil {
				log.Printf("Error sending message: %v", err)
			} else {
				if res.Permission {
					n.Outstanding_Reply -= 1
					log.Print("reply received")
				}
			}
		}
	}
	n.peerLock.RUnlock()
}

func (n *P2PNode) Start() {
	// Simulate sending a message to a peer
	for {
		n.Request_Critical = true
		log.Printf("Node 1 requests access \n")
		n.Our_Timestamp = n.Highest_Timestamp + 1
		n.Outstanding_Reply = n.N - 1
		n.AskAllPeers()
		timeout := 0
		for {
			if timeout == 30 {
				//ReleaseCS()
			}
			if n.Outstanding_Reply == 0 {
				//Enter Critical section
				//Do something
				//Exit Critical section
				//Send Reply to all
				break

			}
			time.Sleep(time.Millisecond * 100)
			timeout++

		}
		time.Sleep(time.Second)
	}
}
