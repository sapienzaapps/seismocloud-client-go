package scsclient

import "errors"

const (
	// Battery indicates that the power source is a battery
	Battery PowerSource = "battery"

	// PowerSupply indicates that the power source is the mains current
	PowerSupply PowerSource = "power-supply"
)

// PowerSource is the power source type
type PowerSource string

// IsValid checks the validity for PowerSource values
func (s PowerSource) IsValid() error {
	switch s {
	case Battery, PowerSupply:
		return nil
	}
	return errors.New("invalid PowerSource type")
}
