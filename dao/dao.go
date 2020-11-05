package dao

import (
	"context"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
)

const (
	rwlockPrefix = "/rwlock"
	rlockPrefix  = "r"
	wlockPrefix  = "w"
)

// Dao : interface for Data Access Object
type Dao interface {
	RLock(res string, owner string) (bool, error)
	Lock(res string, owner string) (bool, error)
	Unlock(res string, owner string) (bool, error)
}

type dao struct {
	etcdCli *clientv3.Client
	kvc     clientv3.KV
}

// Config : configurations for initializing dao
type Config struct {
	Endpoints []string
	Timeout   time.Duration
}

// NewDao : create and return new Dao with configuration
func NewDao(cfg Config) (Dao, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: cfg.Timeout,
	})
	if err != nil {
		return nil, err
	}
	return &dao{
		etcdCli: cli,
		kvc:     clientv3.NewKV(cli),
	}, nil
}

func (d *dao) RLock(res string, owner string) (bool, error) {
	r, err := d.kvc.Txn(context.TODO()).
		If(
			clientv3.Compare(
				clientv3.CreateRevision(strings.Join([]string{
					rwlockPrefix, res, "w",
				}, "/")).WithPrefix(), "=", 0,
			),
			clientv3.Compare(
				clientv3.CreateRevision(strings.Join([]string{
					rwlockPrefix, res, "r", owner,
				}, "/")), "=", 0,
			),
		).
		Then(
			clientv3.OpPut(strings.Join([]string{
				rwlockPrefix, res, "r", owner,
			}, "/"), "1"),
		).
		Else().
		Commit()
	if err != nil || r == nil {
		return false, err
	}
	return r.Succeeded, nil
}

func (d *dao) Lock(res string, owner string) (bool, error) {
	r, err := d.kvc.Txn(context.TODO()).
		If(
			clientv3.Compare(
				clientv3.CreateRevision(strings.Join([]string{
					rwlockPrefix, res, "w",
				}, "/")).WithPrefix(), "=", 0,
			),
			clientv3.Compare(
				clientv3.CreateRevision(strings.Join([]string{
					rwlockPrefix, res, "r",
				}, "/")).WithPrefix(), "=", 0,
			),
		).
		Then(
			clientv3.OpPut(strings.Join([]string{
				rwlockPrefix, res, "w", owner,
			}, "/"), "1"),
		).
		Else().
		Commit()
	if err != nil || r == nil {
		return false, err
	}
	return r.Succeeded, nil
}

func (d *dao) Unlock(res string, owner string) (bool, error) {
	r, err := d.kvc.Delete(context.TODO(), strings.Join([]string{
		rwlockPrefix, res, "r", owner,
	}, "/"))
	if err != nil || r == nil {
		return false, err
	}
	if r.Deleted > 0 {
		return true, nil
	}
	r, err = d.kvc.Delete(context.TODO(), strings.Join([]string{
		rwlockPrefix, res, "w", owner,
	}, "/"))
	if r.Deleted > 0 {
		return true, nil
	}
	return false, nil
}
