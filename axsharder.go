package axsharder

import (
	"context"
	"github.com/go-errors/errors"
)

/*
zed (09.03.2024)
*/

var ErrNotStarted = errors.New("AxSharder not started")

type AxSharder[T IShardMessage] struct {
	shardChannels map[int]chan T
	inChannel     chan T
	count         int
	ctx           context.Context
	cancelFn      context.CancelFunc
	started       bool
}

func NewAxSharderInChan[T IShardMessage](ctx context.Context, count int, shardBufferSize int, in chan T) *AxSharder[T] {
	res := &AxSharder[T]{
		inChannel:     in,
		shardChannels: make(map[int]chan T, shardBufferSize),
		count:         count,
	}
	res.ctx, res.cancelFn = context.WithCancel(ctx)
	for i := 0; i < count; i++ {
		res.shardChannels[i] = make(chan T)
	}

	return res
}

func NewAxSharder[T IShardMessage](ctx context.Context, count int, shardBufferSize int) *AxSharder[T] {
	return NewAxSharderInChan(ctx, count, shardBufferSize, make(chan T, shardBufferSize*count))
}

func (s *AxSharder[T]) Send(msg T) error {
	if s.started {
		s.inChannel <- msg
		return nil
	}
	return ErrNotStarted
}

func (s *AxSharder[T]) Stop() {
	if !s.started {
		return
	}
	s.cancelFn()
}

func (s *AxSharder[T]) C(key int) chan T {
	if key < 0 || key >= s.count {
		return nil
	}
	return s.shardChannels[key]
}

func (s *AxSharder[T]) Start() {
	if s.started {
		return
	}
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case msg := <-s.inChannel:
				key := msg.GetShardKey()
				chanKey := key % s.count
				s.shardChannels[chanKey] <- msg
			}
		}
	}()
}

func (s *AxSharder[T]) ShardCount() int {
	return s.count
}

func (s *AxSharder[T]) IsStarted() bool {
	return s.started
}

func (s *AxSharder[T]) Queues() (int, []int) {
	res := make([]int, s.count)
	for k, v := range s.shardChannels {
		res[k] = len(v)
	}
	return len(s.inChannel), res
}
