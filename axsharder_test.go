package axsharder

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/*
 __    _           ___
|  |  |_|_____ ___|_  |
|  |__| |     | .'|  _|
|_____|_|_|_|_|__,|___|
zed (09.03.2024)
*/

type simpleAxSharderMessage struct {
	User  int
	Value int
}

func (s *simpleAxSharderMessage) GetShardKey() int {
	return s.User
}

func TestAxSharder_Chan(t *testing.T) {
	inChan := make(chan *simpleAxSharderMessage, 100)
	ctx := context.Background()
	sharder := NewAxSharderInChan[*simpleAxSharderMessage](ctx, 4, 100, inChan)
	assert.NotNil(t, sharder)
	for i := 0; i < sharder.ShardCount(); i++ {
		go func(k int) {
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-sharder.C(k):
					t.Logf("Shard:%d Recv: User:%d Value:%d", k, msg.User, msg.Value)
				}
			}
		}(i)
	}
	sharder.Start()
	for i := 0; i < 100; i++ {
		inChan <- &simpleAxSharderMessage{User: i % 10, Value: i}
	}
	time.Sleep(time.Millisecond * 100)
	sharder.Stop()
}
