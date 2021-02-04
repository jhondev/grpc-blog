package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // use it for better error detail

	fmt.Println("Configuring grpc server...")
	fmt.Println("Contecting to mongodb server...")
	connectToMongodb()

	address := "localhost:50051"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("ssl/gogen/cert.pem", "ssl/gogen/key.pem")
	if err != nil {
		log.Fatalf("Failed to load certificates %v", err)
	}
	server := grpc.NewServer(grpc.Creds(creds))
	blogpb.RegisterBlogServiceServer(server, &GrpcServer{})

	go func() {
		fmt.Printf("\nStarting server on %v\n", address)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// wait for ctrl c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until signal is received
	<-ch
	fmt.Println("Stopping the server")
	server.Stop()
	fmt.Println("Closing the listener")
	_ = listener.Close()
	fmt.Println("End of program")
}
