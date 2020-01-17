package entities

// Ip list interface
type IPList interface {
	Add(ip IP) (bool, error)
	Delete(ip IP, mode int) (bool, error)
	Has(ip IP) (bool, error)
	Count() (int, error)
}
