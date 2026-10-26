package base

import (
	"context"
	"micro-scrm/internal/events/entity"
)

type IListener interface {
	Handle(data entity.Data)
	SyncHandle(data entity.Data)
}

type BaseEvent struct {
	Name     string
	Listener []IListener
	Ctx      context.Context
}

// 注册监听者
func (b *BaseEvent) RegisterListener(listenerList ...IListener) {
	for _, listener := range listenerList {
		b.Listener = append(b.Listener, listener)
	}
}

// 并发事件
func (b BaseEvent) Dispatch(data entity.Data) {
	for _, item := range b.Listener {
		go item.Handle(data)
	}
}

// 同步事件
func (b BaseEvent) SyncDispatch(data entity.Data) {
	for _, item := range b.Listener {
		item.SyncHandle(data)
	}
}

type BaseListener struct {
	Name string
	// TODO 通知退出，超时机制
	Ctx context.Context
}
