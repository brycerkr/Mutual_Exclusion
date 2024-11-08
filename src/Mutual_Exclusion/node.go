package main

import (
	"Mutual_Exclusion/m/v2/raalgo"
	"log"
	"net"

	"google.golang.org/grpc"
)

//Node needs to keep list of other nodes in network
//Look at map from ChittyChat - map of nodeid to stream

type P2PServer struct {

}

func start() {
	/** 
	
	The concepts of client and server are baked into grpc
	Only clients can make rpc calls to servers

	So every node must contain both a "client" and a "server"

	Start listening on port x (pass in from main)

	Once each "server" is running, clients on other nodes can connect to that port

    Step 1: Initialize the first node
    1.1 : start the first node's server

    ......

    For nodes: start them as goroutines


    */

    


}

func StartServer(port string) {
    lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
    
    grpcServer := grpc.NewServer()

    server := &P2PServer{

    }

    // Register the ChittyChatServer with the gRPC server
    P2Pnetwork.RegisterP2PNetworkServer(grpcServer, server)

    log.Printf("ChittyChat server started at port :50051")
    if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}

func StartClient() {
    // Establish connection to the gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

}
 




	