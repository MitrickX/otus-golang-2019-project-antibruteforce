package ip

// Ip list interface
type List interface {
	// Add IP into list
	Add(ip IP) (bool, error)
	// Delete IP from list
	Delete(ip IP) (bool, error)
	// Has list this IP
	Has(ip IP) (bool, error)
	// How many IPs in list
	Count() (int, error)
}
