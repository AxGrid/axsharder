package axsharder

import "context"

/*
 __    _           ___
|  |  |_|_____ ___|_  |
|  |__| |     | .'|  _|
|_____|_|_|_|_|__,|___|
zed (09.03.2024)
*/

func NewAxSharderWorker[T IShardMessage](ctx context.Context, count int, shardBufferSize int, worker AxShardWorkerFunc[T]) *AxSharder[T] {
	res := NewAxSharder[T](ctx, count, shardBufferSize)
	for i := 0; i < count; i++ {
		go func(key int) {
			select {
			case <-res.ctx.Done():
				return
			case msg := <-res.C(key):
				worker(res.ctx, msg)
			}
		}(i)
	}
	return res
}

func NewAxSharderChanWorker[T IShardMessage](ctx context.Context, count int, shardBufferSize int, worker AxShardWorkerChanFunc[T]) *AxSharder[T] {
	res := NewAxSharder[T](ctx, count, shardBufferSize)
	for i := 0; i < count; i++ {
		go worker(res.ctx, res.C(i))
	}
	return res
}

func NewAxSharderWorkerInChan[T IShardMessage](ctx context.Context, count int, shardBufferSize int, worker AxShardWorkerFunc[T], in chan T) *AxSharder[T] {
	res := NewAxSharderInChan[T](ctx, count, shardBufferSize, in)
	for i := 0; i < count; i++ {
		go func(key int) {
			select {
			case <-res.ctx.Done():
				return
			case msg := <-res.C(key):
				worker(res.ctx, msg)
			}
		}(i)
	}
	return res
}

func NewAxSharderChanWorkerInChan[T IShardMessage](ctx context.Context, count int, shardBufferSize int, worker AxShardWorkerChanFunc[T], in chan T) *AxSharder[T] {
	res := NewAxSharderInChan[T](ctx, count, shardBufferSize, in)
	for i := 0; i < count; i++ {
		go worker(res.ctx, res.C(i))
	}
	return res
}
