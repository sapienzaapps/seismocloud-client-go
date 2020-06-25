package scsclient

import "errors"

const (
	Battery     PowerSource = "battery"
	PowerSupply PowerSource = "power-supply"
)

type PowerSource string

func (s PowerSource) IsValid() error {
	switch s {
	case Battery, PowerSupply:
		return nil
	}
	return errors.New("invalid PowerSource type")
}
