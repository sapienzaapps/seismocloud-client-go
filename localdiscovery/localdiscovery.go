package localdiscovery

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"strings"
)

const (
	PKTTYPE_DISCOVERY       = 1
	PKTTYPE_DISCOVERY_REPLY = 2
)

type Sensor struct {
	DeviceID string
	Model    string
	Version  string
	RemoteIP net.IP
}

func Discovery(ctx context.Context, sensors chan<- Sensor) ([]Sensor, error) {
	var ret []Sensor

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 62001,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		return ret, errors.Wrap(err, "creating UDP listener")
	}
	defer func() {
		_ = conn.Close()
	}()

	msgbuffer := make([]byte, 5+1+1)
	copy(msgbuffer, "INGV\000")
	msgbuffer[5] = PKTTYPE_DISCOVERY

	broadcastDestination := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 62001}
	_, err = conn.WriteToUDP(msgbuffer, &broadcastDestination)
	if err != nil {
		return ret, errors.Wrap(err, "sending UDP discovery packet")
	}

	go func() {
		var buf [1024]byte
		for {
			rlen, remote, err := conn.ReadFromUDP(buf[:])
			if err != nil {
				// TODO: handle this error
				return
				//return ret, errors.Wrap(err, "reading from UDP")
			}

			if rlen > 5 && bytes.Compare(buf[:5], []byte("INGV\000")) == 0 {
				switch buf[5] {
				case PKTTYPE_DISCOVERY_REPLY:
					dev := Sensor{
						DeviceID: fmt.Sprintf("%x", buf[6:6+6]),
						Model:    strings.Trim(string(buf[6+6+4:6+6+4+8]), "\000"),
						Version:  string(buf[6+6 : 6+6+4]),
						RemoteIP: remote.IP,
					}
					ret = append(ret, dev)

					if sensors != nil {
						sensors <- dev
					}
				}
			}
		}
	}()

	<-ctx.Done()
	return ret, nil
}
