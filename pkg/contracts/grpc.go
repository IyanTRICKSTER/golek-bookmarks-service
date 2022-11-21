package contracts

import (
	"context"
	"golek_bookmark_service/pkg/models"
	ps "golek_bookmark_service/pkg/models/proto_schema"
)

type GRPCClient interface {
	Dial() (ps.PostServiceClient, error)
}

type GRPCPostService interface {
	Fetch(ctx context.Context, postIDs []string) ([]models.Post, error)
	GRPCClient
}
