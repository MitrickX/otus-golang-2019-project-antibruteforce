package ip

import (
	"context"

	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

type List struct {
	list []entities.IP
}

func NewList() *List {
	return &List{}
}

func (l *List) Add(ctx context.Context, ip entities.IP) error {
	index := l.search(ip)
	if index >= 0 {
		return nil
	}
	l.list = append(l.list, ip)
	return nil
}

func (l *List) Delete(ctx context.Context, ip entities.IP) error {
	index := l.search(ip)
	if index < 0 {
		return nil
	}
	l.list[index] = l.list[len(l.list)-1]
	l.list = l.list[:len(l.list)-1]
	return nil
}

func (l *List) Has(ctx context.Context, ip entities.IP) (bool, error) {
	index := l.search(ip)
	if index >= 0 {
		return true, nil
	}
	return false, nil
}

func (l *List) IsConform(ctx context.Context, ip entities.IP) (bool, error) {
	ip = ip.DropMaskPart()
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

func (l *List) search(ip entities.IP) int {
	for index, value := range l.list {
		if value == ip {
			return index
		}
	}
	return -1
}

func (l *List) Count(context.Context) (int, error) {
	return len(l.list), nil
}
