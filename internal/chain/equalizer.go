package chain

import (
	"context"
	"github.com/ethereum/go-ethereum/common/math"
	"token-payment/internal/dao"
)

const (
	scoreMax = math.MaxInt64 - 1000
	skipStep = 1000
)

type Equalizer struct {
	Name    string
	Members []string
}

func NewEqualizer(ctx context.Context, name string, members []string) *Equalizer {
	e := &Equalizer{
		Name:    name,
		Members: members,
	}
	e.Init(ctx)
	return e
}

func (e *Equalizer) Init(ctx context.Context) {
	for _, member := range e.Members {
		dao.Redis.ZIncrBy(ctx, e.Name, 0, member)
	}
}

func (e *Equalizer) reset(ctx context.Context) {
	// 将score全部设置为0
	dao.Redis.Del(ctx, e.Name)
	e.Init(ctx)
}

func (e *Equalizer) Get(ctx context.Context) (string, error) {
	// 从有序集合中获取最小的一个
	members, err := dao.Redis.ZRangeWithScores(ctx, e.Name, 0, 0).Result()
	if err != nil || len(members) == 0 {
		return "", err
	}
	m := members[0]
	if m.Score > scoreMax {
		e.reset(ctx)
	} else {
		dao.Redis.ZIncrBy(ctx, e.Name, 1, m.Member.(string))
	}
	return m.Member.(string), nil
}

func (e *Equalizer) Skip(ctx context.Context, member string) {
	dao.Redis.ZIncrBy(context.Background(), e.Name, skipStep, member)
}
