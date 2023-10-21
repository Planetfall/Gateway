// Package search contains the music researcher controller.
// It aims to interact using gRPC with the music researcher service.
package search

import (
	"fmt"

	"github.com/planetfall/gateway/internal/connection/grpc"
	"github.com/planetfall/gateway/internal/controller"
	pb "github.com/planetfall/genproto/pkg/musicresearcher/v1"
)

// SearchController is used to interact with the music researcher service.
type SearchController struct {
	// Reference to the base controller type
	controller.Controller

	// The connection to setup the client and authenticate the request context
	conn *grpc.Connection

	// The generate protobuf client for the
	// [github.com/planetfall/musicresearcher] service
	client pb.MusicResearcherClient
}

// SearchControllerOptions holds the parameters for the ResearcherController
// builder
type SearchControllerOptions struct {
	// The [controller] builder parameters
	ControllerOptions controller.ControllerOptions

	// Insecure for [grpc] connection builder parameters
	Insecure bool
}

// NewSearchController buids a new MusicResearcher controller.
// It setup the GRPC connection and the protobuf client.
func NewSearchController(
	opt SearchControllerOptions) (*SearchController, error) {

	// initialize the base type
	ctrl := controller.NewController(opt.ControllerOptions)

	// setup the connection
	conn, err := grpc.NewConnection(grpc.ConnectionOptions{
		Target:   opt.ControllerOptions.Target,
		Insecure: opt.Insecure,
	})
	if err != nil {
		return nil, fmt.Errorf("connection.NewConnection: %v", err)
	}

	// setup the client
	client := pb.NewMusicResearcherClient(conn.GrpcConn())

	return &SearchController{
		Controller: ctrl,
		client:     client,
		conn:       conn,
	}, nil
}

// Close terminates the inner GRPC connection
func (c *SearchController) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("connection.Close: %v", err)
	}
	return nil
}
