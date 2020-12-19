package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-blog/blogpb"
	"log"
)

func main() {
	fmt.Println("Connecting to grpc server...")
	adderess := "localhost:50051"
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial(adderess, opts)
	if err != nil {
		log.Fatalf("\nError connecting to grpc server %v\n", adderess)
	}
	defer conn.Close()

	client := blogpb.NewBlogServiceClient(conn)
	fmt.Println("Connected")
	blogReq := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Title:    "Creating a blog grpc service",
			Content:  "This is the content for the blog",
			AuthorId: "jhondev",
		},
	}
	fmt.Println("Creating a blog...")
	response, err := client.CreateBlog(context.Background(), blogReq)
	if err != nil {
		log.Fatalf("\nError creating the blog: %v\n", err)
	}
	log.Printf("\nBlog created: %v\n", response)
}
