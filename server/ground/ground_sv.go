package ground

import (
	"context"
	"feimumoke/farming/v2/api/server"
	"github.com/google/uuid"
	"math/rand"
)

type GroundSv struct {
}

func (g GroundSv) SelectGround(ctx context.Context, req *server.GroundReq) (*server.GroundRsp, error) {
	return &server.GroundRsp{
		GroundId: uuid.New().String(),
		Position: uuid.New().String(),
		Owner:    req.Owner,
		Price:    rand.Int63n(100),
	}, nil
}

func (g GroundSv) mustEmbedUnimplementedGroundServiceServer() {
	panic("implement me")
}
