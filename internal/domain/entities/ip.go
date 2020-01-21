package entities

import (
	"context"
	"net"
	"strings"
)

type IP string

func (ip IP) DropMaskPart() IP {
	parts := strings.Split(string(ip), "/")
	return IP(parts[0])
}

func (ip IP) HasMaskPart() bool {
	parts := strings.Split(string(ip), "/")
	return len(parts) > 1
}

func (ip IP) ParseAsCIDR() (net.IP, *net.IPNet, error) {
	return net.ParseCIDR(string(ip))
}

func (ip IP) Parse() net.IP {
	return net.ParseIP(string(ip.DropMaskPart()))
}

func (ip IP) IsConform(checkedIP IP) bool {
	if !ip.HasMaskPart() {
		return false
	}

	if checkedIP.HasMaskPart() {
		return false
	}

	_, mask, err := ip.ParseAsCIDR()
	if err != nil {
		return false
	}

	checkedNetIP := checkedIP.Parse()
	if checkedNetIP == nil {
		return false
	}

	return mask.Contains(checkedNetIP)
}

// Ip list interface
type IPList interface {
	// Add IP into list
	Add(ctx context.Context, ip IP) error
	// Delete IP from list
	Delete(ctx context.Context, ip IP) error
	// Has list this IP
	Has(ctx context.Context, ip IP) (bool, error)
	// Is IP conformed list. More broader concept than Has - checking matching among subnet IPs in list
	IsConform(ctx context.Context, ip IP) (bool, error)
	// How many IPs in list
	Count(ctx context.Context) (int, error)
}
