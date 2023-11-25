package util

// constants for all supported units
const (
	KWH = "kWh"
)

// IsSupportedUnit returns true if the Unit is supported
func IsSupportedUnit(unit string) bool {
	switch unit {
	case KWH:
		return true
	}
	return false
}
