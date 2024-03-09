package axsharder

import "context"

/*
 __    _           ___
|  |  |_|_____ ___|_  |
|  |__| |     | .'|  _|
|_____|_|_|_|_|__,|___|
zed (09.03.2024)
*/

type IShardMessage interface {
	GetShardKey() int
}

type AxShardWorkerChanFunc[T IShardMessage] func(context.Context, chan T)
type AxShardWorkerFunc[T IShardMessage] func(context.Context, T)


type AxSharderOpts struct {
	ShardCount int
	ShardOutBufferSize int
}