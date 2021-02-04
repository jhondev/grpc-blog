package main

import (
	"context"
	"fmt"
	"grpc-blog/blogpb"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct{}

func (server GrpcServer) ReadBlog(_ context.Context, request *blogpb.ReadBlogRequest) (*blogpb.ReadBlogResponse, error) {
	fmt.Println("\nReading blog request")
	blogId := request.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse id %s: %s", blogId, err)
	}
	blog := &Blog{}
	filter := bson.D{{"_id", oid}}
	if err := blogs.FindOne(context.Background(), filter).Decode(blog); err != nil {
		fmt.Printf("Reading blog error: %v", err)
		return nil, status.Errorf(codes.NotFound, "Blog not found (_id: %s)", oid)
	}

	fmt.Println("Blog read")
	return &blogpb.ReadBlogResponse{
		Blog: &blogpb.Blog{
			Id:       blogId,
			AuthorId: blog.AuthorID,
			Content:  blog.Content,
			Title:    blog.Title,
		}}, nil
}

func (server GrpcServer) CreateBlog(_ context.Context, request *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	fmt.Println("\nCreating a blog request")
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

func (server GrpcServer) UpdateBlog(_ context.Context, request *blogpb.UpdateBlogRequest) (*blogpb.UpdateBlogResponse, error) {
	reqBlog := request.GetBlog()
	fmt.Printf("\nUpdating a blog:(id:%v)\n", reqBlog.Id)
	oid, err := primitive.ObjectIDFromHex(reqBlog.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse id %s: %s", reqBlog.Id, err)
	}

	blog := &Blog{
		ID:       oid,
		AuthorID: reqBlog.AuthorId,
		Content:  reqBlog.Content,
		Title:    reqBlog.Title,
	}
	filter := bson.D{{"_id", oid}}
	result, err := blogs.ReplaceOne(context.Background(), filter, blog)
	if err != nil {
		fmt.Printf("Updating blog error: %v", err)
	}
	if result.ModifiedCount == 0 {
		return nil, status.Errorf(codes.Internal, "Couldn't update blog (_id: %s)", oid)
	}
	fmt.Println("Blog updated")
	return &blogpb.UpdateBlogResponse{BlogId: reqBlog.Id}, nil
}

func (server GrpcServer) DeleteBlog(_ context.Context, request *blogpb.DeleteBlogRequest) (*emptypb.Empty, error) {
	fmt.Println("\nDeleting blog request")
	blogId := request.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Cannot parse id %s: %s", blogId, err)
	}
	filter := bson.D{{"_id", oid}}
	_, err = blogs.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error deleting blog (id:%s): %s", blogId, err)
	}
	fmt.Println("Blog deleted")
	return &emptypb.Empty{}, nil
}
