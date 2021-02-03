package main

import (
	"context"
	"fmt"
	"grpc-blog/blogpb"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type Blog struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

var blogs *mongo.Collection

type grpcServer struct{}

func (server grpcServer) CreateBlog(ctx context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("Creating a blog...")
	blog := request.GetBlog()
	data := &Blog{
		AuthorID: blog.AuthorId,
		Content:  blog.Content,
		Title:    blog.Title,
	}
	result, err := blogs.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "\nError creating a new blog: %v\n", err)
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Error casting objectid")
	}
	fmt.Println("Blog created")
	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			Id:       oid.Hex(),
			AuthorId: blog.AuthorId,
			Content:  blog.Content,
			Title:    blog.Title,
		}}, nil
}

func connectToMongodb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("\nError connecting to mongodb server\n")
	}
	blogs = client.Database("blogdb").Collection("blogs")
}

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
	blogpb.RegisterBlogServiceServer(server, &grpcServer{})

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
