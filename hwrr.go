package hwrr

// Backend represents a server in the load balancer.
type Backend struct {
	Name   string
	Weight int
}

// HWRR represents the load balancer.
type HWRR struct {
	backends []*Backend
}

// NewHWRR creates a new load balancer with provided backends.
func NewHWRR(backends []*Backend) *HWRR {
	hwrr := &HWRR{}

	for _, b := range backends {
		if b.Weight <= 0 {
			continue
		}

		for i := 0; i < b.Weight; i++ {
			hwrr.backends = append(hwrr.backends, b)
		}
	}

	return hwrr
}

// GetNextBackend returns the next backend server using weighted round-robin.
func (hwrr *HWRR) GetNextBackend(hash int) *Backend {
	if len(hwrr.backends) == 0 {
		return nil
	}

	index := hash % len(hwrr.backends)

	selected := hwrr.backends[index]

	return selected
}
