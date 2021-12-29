package localdiscovery

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
)

// Beacon will start a goroutine that sends a reply to each discovery, using the legacy protocol.
// Note that the deviceId should be a hex value of 12 character, model a string of 8 and version a string of 4 chars
func Beacon(ctx context.Context, deviceID string, model string, version string) error {
	if len(deviceID) != 12 {
		return errors.New("device ID format error")
	}

	deviceIDValue, err := hex.DecodeString(deviceID)
	if err != nil {
		return fmt.Errorf("device ID format error: %w", err)
	}

	if len(model) > 8 || len(model) == 0 {
		return errors.New("model size wrong")
	}
	if len(version) > 4 || len(version) == 0 {
		return errors.New("version size wrong")
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 62001,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		return fmt.Errorf("creating UDP listener: %w", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	var errchannel = make(chan error, 1)
	go func() {
		var buf [1024]byte
		for {
			rlen, remote, err := conn.ReadFromUDP(buf[:])
			if err != nil {
				errchannel <- fmt.Errorf("error reading from UDP socket: %w", err)
				return
			}

			if rlen > 5 && bytes.Equal(buf[:5], []byte("INGV\000")) && buf[5] == PacketDiscovery {
				msgbuffer := make([]byte, 5+1+6+4+8)
				copy(msgbuffer, "INGV\000")
				msgbuffer[5] = PacketDiscoveryReply
				copyByteN(msgbuffer, 5+1, deviceIDValue, 6)
				copyByteN(msgbuffer, 5+1+6, []byte(version), 4)
				copyByteN(msgbuffer, 5+1+6+4, []byte(model), 8)
				_, err = conn.WriteToUDP(msgbuffer, remote)
				if err != nil {
					errchannel <- fmt.Errorf("error sending reply to UDP socket: %w", err)
					return
				}
			}
		}
	}()

	select {
	case err := <-errchannel:
		return err
	case <-ctx.Done():
	}
	return nil
}

func copyByteN(dst []byte, dstoffset int, src []byte, max int) {
	for idx, r := range src {
		if idx >= max {
			return
		}
		dst[dstoffset+idx] = r
	}
}
