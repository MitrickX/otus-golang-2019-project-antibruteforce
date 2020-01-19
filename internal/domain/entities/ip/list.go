package ip

import "context"

// Ip list interface
type List interface {
	// Add IP into list
	Add(ctx context.Context, ip IP) error
	// Delete IP from list
	Delete(ctx context.Context, ip IP) error
	// Has list this IP
	Has(ctx context.Context, ip IP) (bool, error)
	// How many IPs in list
	Count(ctx context.Context) (int, error)
}
