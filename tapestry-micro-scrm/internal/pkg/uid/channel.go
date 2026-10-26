package uid

import (
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/threading"
	"sync"
	"time"
)

type Channel struct {
	sync.RWMutex
	ids chan int64
}

var ChannelId *Channel

func NewChannel() *Channel {
	ChannelId = &Channel{
		ids: make(chan int64, 1000),
	}
	return ChannelId
}

// GenChannelId 异步生产全局id
func (c *Channel) GenChannelId() {
	threading.GoSafe(func() error {
		for {
			if len(c.ids) < 300 {
				for _, id := range Ids(ctx.New(), 1000) {
					if len(c.ids) == 1000 {
						break
					}
					c.ids <- id
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	})
}

func (c *Channel) GetChannelId() int64 {
	c.Lock()
	defer c.Unlock()
	id := <-c.ids
	if id == 0 {
		id = ID(ctx.New())
	}
	return id
}
