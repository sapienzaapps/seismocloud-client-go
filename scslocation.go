package scsclient

import "github.com/shopspring/decimal"

type SCSLocation struct {
	Lat decimal.Decimal
	Lng decimal.Decimal
}

func (loc *SCSLocation) IsValid() bool {
	return !(loc.Lat.IsZero() && loc.Lng.IsZero())
}
