package ip

import (
	"context"
	"fmt"
	"sync"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

type List struct {
	list []entities.IP
	mx   sync.RWMutex
}

func NewList() *List {
	return &List{
		mx: sync.RWMutex{},
	}
}

func (l *List) Add(ctx context.Context, ip entities.IP) error {
	l.mx.Lock()
	defer l.mx.Unlock()

	index := l.search(ip)
	if index >= 0 {
		return nil
	}

	l.list = append(l.list, ip)

	return nil
}

func (l *List) Delete(ctx context.Context, ip entities.IP) error {
	l.mx.Lock()
	defer l.mx.Unlock()

	index := l.search(ip)
	if index < 0 {
		return nil
	}

	l.list[index] = l.list[len(l.list)-1]
	l.list = l.list[:len(l.list)-1]

	return nil
}

func (l *List) Has(ctx context.Context, ip entities.IP) (bool, error) {
	l.mx.RLock()
	defer l.mx.RUnlock()

	index := l.search(ip)
	if index >= 0 {
		return true, nil
	}

	return false, nil
}

func (l *List) IsConform(ctx context.Context, ip entities.IP) (bool, error) {
	l.mx.RLock()
	defer l.mx.RUnlock()

	if ip.HasMaskPart() {
		return false, fmt.Errorf("expected pure IPBits (without mask) instread of `%s`", ip)
	}

	for _, ipInList := range l.list {
		if ipInList == ip {
			return true, nil
		}

		if ipInList.IsConform(ip) {
			return true, nil
		}
	}

	return false, nil
}

func (l *List) Count(context.Context) (int, error) {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return len(l.list), nil
}

func (l *List) Clear(context.Context) error {
	l.mx.Lock()
	defer l.mx.Unlock()
	l.list = l.list[0:0]

	return nil
}

func (l *List) search(ip entities.IP) int {
	for index, value := range l.list {
		if value == ip {
			return index
		}
	}

	return -1
}
