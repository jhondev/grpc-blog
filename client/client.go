package main

import (
	"context"
	"fmt"
	"grpc-blog/blogpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("\nConnecting to grpc server...")
	adderess := "localhost:50051"
	creds, err := credentials.NewClientTLSFromFile("ssl/gogen/cert.pem", "")
	if err != nil {
		log.Fatalf("\nError loading certificates %v\n", adderess)
	}
	// opts := grpc.WithInsecure()
	opts := grpc.WithTransportCredentials(creds)
	conn, err := grpc.Dial(adderess, opts)
	if err != nil {
		log.Fatalf("\nError connecting to grpc server %v\n", adderess)
	}
	defer conn.Close()

	client := blogpb.NewBlogServiceClient(conn)
	fmt.Println("Connected")

	result := CreateBlog(client)
	blog := ReadBlog(client, result.Blog.Id)
	UpdateBlog(client, blog)
	ReadBlog(client, result.Blog.Id)
	DeleteBlog(client, result.Blog.Id)
	ReadBlog(client, result.Blog.Id)
}

func CreateBlog(client blogpb.BlogServiceClient) *blogpb.CreateBlogResponse {
	blogReq := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			Title:    "Creating a blog grpc service",
			Content:  "This is the content for the blog",
			AuthorId: "jhondev",
		},
	}
	fmt.Println("\nCreating a blog...")
	result, err := client.CreateBlog(context.Background(), blogReq)
	if err != nil {
		log.Fatalf("\nError creating the blog: %v\n", err)
	}
	log.Printf("\nBlog created: %v\n", result)

	return result
}

func ReadBlog(client blogpb.BlogServiceClient, blogId string) *blogpb.Blog {
	fmt.Println("\nReading a blog...")
	result, err := client.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: blogId})
	if err != nil {
		log.Fatalf("\nError reading blog: %s\n", err)
	}
	log.Printf("\nBlog info: %v\n", result)
	return result.Blog
}

func UpdateBlog(client blogpb.BlogServiceClient, blog *blogpb.Blog) string {
	fmt.Println("\nUpdating a blog...")
	blog.Title = "Title edited"
	result, err := client.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("\nError updating blog: %s\n", err)
	}
	log.Printf("Blog updated (id:%s)", result.BlogId)
	return result.BlogId
}

func DeleteBlog(client blogpb.BlogServiceClient, blogId string) {
	fmt.Println("\nDeleting a blog...")
	_, err := client.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogId})
	if err != nil {
		log.Fatalf("\nError deleting blog: %s\n", err)
	}
	log.Printf("Blog deleted (id:%s)", blogId)
}
