package memory

import (
	"github.com/mitrickx/otus-golang-2019-project-antibruteforce/internal/domain/entities"
)

type IPList struct {
	list []entities.IP
}

func NewIPList() *IPList {
	return &IPList{}
}

func (l *IPList) Add(ip entities.IP) (bool, error) {
	index := l.search(ip)
	if index >= 0 {
		return true, nil
	}
	l.list = append(l.list, ip)
	return true, nil
}

func (l *IPList) Delete(ip entities.IP) (bool, error) {
	index := l.search(ip)
	if index < 0 {
		return true, nil
	}
	l.list[index] = l.list[len(l.list)-1]
	l.list = l.list[:len(l.list)-1]
	return true, nil
}

func (l *IPList) Has(ip entities.IP) (bool, error) {
	index := l.search(ip)
	if index >= 0 {
		return true, nil
	}
	return false, nil
}

func (l *IPList) search(ip entities.IP) int {
	for index, value := range l.list {
		if value == ip {
			return index
		}
	}
	return -1
}

func (l *IPList) Count() (int, error) {
	return len(l.list), nil
}
