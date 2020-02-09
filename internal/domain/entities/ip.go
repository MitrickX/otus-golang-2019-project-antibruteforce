package entities

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IP type
type IP string

// New is general constructor, support ip with mask part or without
func New(ip string) (IP, error) {
	if net.ParseIP(ip) != nil {
		return IP(ip), nil
	}

	if !IP(ip).HasMaskPart() {
		return "", fmt.Errorf("invalid ip `%s`", ip)
	}

	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		return "", fmt.Errorf("invalid ip `%s`: %w", ip, err)
	}

	return IP(ip), nil
}

// NewWithMaskPart is a constructor for ip with mask part, e.g. 127.0.0.0/24
func NewWithMaskPart(ip string) (IP, error) {
	if !IP(ip).HasMaskPart() {
		return "", fmt.Errorf("invalid ip `%s`, mask part is required", ip)
	}

	_, _, err := net.ParseCIDR(ip)
	if err != nil {
		return "", fmt.Errorf("invalid ip `%s`: %w", ip, err)
	}

	return IP(ip), nil
}

// NewWithoutMaskPart is a constructor for ip with mask part, e.g. 127.0.0.1
func NewWithoutMaskPart(ip string) (IP, error) {
	if IP(ip).HasMaskPart() {
		return "", fmt.Errorf("invalid ip `%s`, mask part is not allowed", ip)
	}

	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid ip `%s`", ip)
	}

	return IP(ip), nil
}

// DropMaskPart drops mask part from IP
func (ip IP) DropMaskPart() IP {
	parts := strings.Split(string(ip), "/")
	return IP(parts[0])
}

// HasMaskPart checks is IP has mask
func (ip IP) HasMaskPart() bool {
	parts := strings.Split(string(ip), "/")
	//nolint:gomnd
	return len(parts) > 1
}

// GetMaskPart returns mask part of IP
func (ip IP) GetMaskPart() string {
	parts := strings.Split(string(ip), "/")
	//nolint:gomnd
	if len(parts) > 1 {
		return parts[1]
	}

	return ""
}

// GetMaskAsInt returns mark part and try to convert it to int
func (ip IP) GetMaskAsInt() (int, error) {
	mask := ip.GetMaskPart()
	return strconv.Atoi(mask)
}

// ParseAsCIDR parses CIDR IP (IP with mask)
func (ip IP) ParseAsCIDR() (net.IP, *net.IPNet, error) {
	return net.ParseCIDR(string(ip))
}

// Parse parses IP and return net.IP
func (ip IP) Parse() net.IP {
	return net.ParseIP(string(ip.DropMaskPart()))
}

// IsConform checks is IP contains checked IP, if current IP is not subnet IP method return false
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

// IPList interface for data structure that represent list of IPs
type IPList interface {
	// Add IP into list
	Add(ctx context.Context, ip IP) error
	// Delete IP from list
	Delete(ctx context.Context, ip IP) error
	// Has list this IP
	Has(ctx context.Context, ip IP) (bool, error)
	// IsConform check is IP conformed list. More broader concept than Has - checking matching among subnet IPs in list
	IsConform(ctx context.Context, ip IP) (bool, error)
	// Count how many IPs in list
	Count(ctx context.Context) (int, error)
	// Clear all IPs in list
	Clear(ctx context.Context) error
}
