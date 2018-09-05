package scsclient

type SCSLocation struct {
	Lat float64
	Lng float64
}

func (loc *SCSLocation) IsValid() bool {
	return loc.Lat != 0 || loc.Lng != 0
}
