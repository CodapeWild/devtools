package session

import (
	"context"
	"sync"
	"time"
)

type item struct {
	value    interface{}
	canceler context.CancelFunc
}

type MapStore struct {
	*sync.RWMutex
	data map[string]*item
}

func NewMapStore() *MapStore {
	return &MapStore{&sync.RWMutex{}, make(map[string]*item)}
}

func (this *MapStore) Store(token string, value interface{}, expsec int64) (err error) {
	this.Lock()
	defer this.Unlock()

	ctx, canceler := context.WithTimeout(context.Background(), time.Duration(expsec)*time.Second)
	this.data[token] = &item{value: value, canceler: canceler}

	go func(ctx context.Context, token string) {
		select {
		case <-ctx.Done():
			this.Lock()
			defer this.Unlock()

			delete(this.data, token)
		}
	}(ctx, token)

	return nil
}

func (this *MapStore) Retrieve(token string) (value interface{}, err error) {
	return this.data[token].value, nil
}

func (this *MapStore) Have(token string) bool {
	return this.data[token] != nil
}

func (this *MapStore) Remove(token string) {
	if this.data[token] != nil {
		this.data[token].canceler()
	}
}
