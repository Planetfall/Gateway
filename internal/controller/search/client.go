package search

import (
	"context"

	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
	grpc "google.golang.org/grpc"
)

// The client used by the controller to interfact with the microservice
type Client interface {
	Search(context.Context, *pb.Parameters,
		...grpc.CallOption) (*pb.Results, error)

	GetGenreList(context.Context, *pb.Empty,
		...grpc.CallOption) (*pb.GenreList, error)
}
