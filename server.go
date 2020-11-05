package main

import (
	"context"
	"time"

	"github.com/MrCroxx/etclock/dao"

	pb "github.com/MrCroxx/etclock/proto"
)

const (
	rwlockPrefix = "rwlock"
)

// type LockerServer interface {
// 	Lock(context.Context, *LockRequest) (*LockReply, error)
// 	RLock(context.Context, *LockRequest) (*LockReply, error)
// 	Unlock(context.Context, *LockRequest) (*LockReply, error)
// }

type glock struct {
	dao dao.Dao
}

// LockerServerConfig
type LockerServerConfig struct {
	Endpoints []string
}

func NewLockerServer(cfg LockerServerConfig) (pb.LockerServer, error) {
	d, err := dao.NewDao(dao.Config{
		Endpoints: cfg.Endpoints,
		Timeout:   time.Second * 3,
	})
	if err != nil {
		return nil, err
	}
	return &glock{
		dao: d,
	}, nil
}

func (g *glock) Lock(ctx context.Context, req *pb.LockRequest) (rpy *pb.LockReply, err error) {
	ok, err := g.dao.Lock(req.GetResource(), req.GetRequester())
	return &pb.LockReply{
		Ok: ok,
	}, err
}

func (g *glock) RLock(ctx context.Context, req *pb.LockRequest) (rpy *pb.LockReply, err error) {
	ok, err := g.dao.RLock(req.GetResource(), req.GetRequester())
	return &pb.LockReply{
		Ok: ok,
	}, err
}

func (g *glock) Unlock(ctx context.Context, req *pb.LockRequest) (rpy *pb.LockReply, err error) {
	ok, err := g.dao.Unlock(req.GetResource(), req.GetRequester())
	return &pb.LockReply{
		Ok: ok,
	}, err
}
