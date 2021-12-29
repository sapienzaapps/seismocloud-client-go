package localdiscovery

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strings"
)

const (
	PacketDiscovery      = 1
	PacketDiscoveryReply = 2
)

type Sensor struct {
	DeviceID string
	Model    string
	Version  string
	RemoteIP net.IP
}

func Scan(ctx context.Context, sensors chan<- Sensor) ([]Sensor, error) {
	var ret []Sensor

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 62001,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		return ret, fmt.Errorf("creating UDP listener: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	msgbuffer := make([]byte, 5+1+1)
	copy(msgbuffer, "INGV\000")
	msgbuffer[5] = PacketDiscovery

	broadcastDestination := net.UDPAddr{IP: net.IPv4(255, 255, 255, 255), Port: 62001}
	_, err = conn.WriteToUDP(msgbuffer, &broadcastDestination)
	if err != nil {
		return ret, fmt.Errorf("sending UDP discovery packet: %w", err)
	}

	var errchannel = make(chan error, 1)
	go func() {
		var buf [1024]byte
		for {
			rlen, remote, err := conn.ReadFromUDP(buf[:])
			if err != nil {
				errchannel <- fmt.Errorf("error reading from UDP socket: %w", err)
				return
			}

			if rlen > 5 && bytes.Equal(buf[:5], []byte("INGV\000")) && buf[5] == PacketDiscoveryReply {
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
	}()

	select {
	case err := <-errchannel:
		return ret, err
	case <-ctx.Done():
	}
	return ret, nil
}
