package ground

import (
	"context"
	"feimumoke/farming/v2/api/server"
	"feimumoke/farming/v2/framework/util"
	"github.com/google/uuid"
	"math/rand"
)

func NewGroundSv() *GroundSv {
	return &GroundSv{}
}

type GroundSv struct {
	*server.UnimplementedGroundServiceServer
}

func (g GroundSv) SelectGround(ctx context.Context, req *server.GroundReq) (*server.GroundRsp, error) {
	return &server.GroundRsp{
		GroundId: uuid.New().String(),
		Position: util.GetHostString(),
		Owner:    req.Owner,
		Price:    rand.Int63n(100),
	}, nil
}
