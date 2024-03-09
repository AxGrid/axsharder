Ax Sharder
==========

Ax Sharder is a simple tools to shard channel messages into multiple channels by *ShardKey*.

## Usage

```go
import "github.com/axgrid/axsharder"

type simpleAxSharderMessage struct {
    User  int
    Value int
}

func (s *simpleAxSharderMessage) GetShardKey() int {
    return s.User
}

func main() {
    inChan := make(chan *simpleAxSharderMessage, 100)
    ctx := context.Background()
    sharder := axsharder.NewAxSharderInChan[*simpleAxSharderMessage](ctx, 4, 100, inChan)
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
```