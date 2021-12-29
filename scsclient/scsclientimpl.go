package scsclient

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const clientTimeout = 1500 * time.Millisecond

type _clientimpl struct {
	opts            ClientOptions
	mqttc           mqtt.Client
	aliveTicker     *time.Ticker
	aliveTickerStop chan int
}

func (c *_clientimpl) IsConnected() bool {
	return c.mqttc.IsConnected()
}

func (c *_clientimpl) Close() error {
	c.aliveTickerStop <- 0

	c.mqttc.Publish(fmt.Sprintf("sensor/%s/disconnect", c.opts.DeviceID), 0, false, "y").WaitTimeout(clientTimeout)
	c.mqttc.Disconnect(1000)
	return nil
}

func (c *_clientimpl) Connect() error {
	conntoken := c.mqttc.Connect()
	commandOk := conntoken.WaitTimeout(clientTimeout)
	if commandOk && conntoken.Error() != nil {
		return conntoken.Error()
	} else if !commandOk {
		return errors.New("connect timeout")
	}

	// Register internal ticker for Alive
	c.aliveTicker = time.NewTicker(14 * time.Minute)
	go func() {
		for {
			select {
			case <-c.aliveTicker.C:
				err := c.SendAlive()
				if err != nil {
					c.opts.Logger.WithError(err).Error("Error sending ALIVE command")
				}
			case <-c.aliveTickerStop:
				c.aliveTicker.Stop()
				// Avoid goroutine leak on stop
				return
			}
		}
	}()

	// Register MQTT handlers
	c.mqttc.SubscribeMultiple(map[string]byte{
		fmt.Sprintf("sensor/%s/sigma", c.opts.DeviceID.String()): byte(2),
		fmt.Sprintf("sensor/%s/reboot", c.opts.DeviceID):         byte(2),
		fmt.Sprintf("sensor/%s/timesync", c.opts.DeviceID):       byte(2),
		fmt.Sprintf("sensor/%s/stream", c.opts.DeviceID):         byte(2),
		fmt.Sprintf("sensor/%s/probespeed", c.opts.DeviceID):     byte(2),
	}, func(client mqtt.Client, message mqtt.Message) {
		if len(message.Payload()) == 0 {
			return
		}

		topicparts := strings.SplitN(message.Topic(), "/", 3)
		command := topicparts[2]
		messageReceivedAt := time.Now()

		switch command {
		case "sigma":
			if c.opts.OnNewSigma != nil {
				sigma, err := strconv.ParseFloat(string(message.Payload()), 64)
				if err != nil {
					c.opts.Logger.WithError(err).Error("Error in SIGMA payload")
				} else {
					c.opts.OnNewSigma(c, sigma)
				}
			}
		case "reboot":
			if c.opts.OnReboot != nil {
				c.opts.OnReboot(c)
			}
		case "timesync":
			if c.opts.OnTimeReceived != nil {
				tsparts := strings.Split(string(message.Payload()), ";")
				if len(tsparts) != 3 {
					c.opts.Logger.Error("Error in TIMESYNC payload: wrong number of arguments")
					return
				}
				t0, err := strconv.ParseInt(tsparts[0], 10, 64)
				if err != nil {
					c.opts.Logger.WithError(err).Error("Error in TIMESYNC T0")
					return
				}
				t1, err := strconv.ParseInt(tsparts[1], 10, 64)
				if err != nil {
					c.opts.Logger.WithError(err).Error("Error in TIMESYNC T1")
					return
				}
				t2, err := strconv.ParseInt(tsparts[2], 10, 64)
				if err != nil {
					c.opts.Logger.WithError(err).Error("Error in TIMESYNC T2")
					return
				}
				t3 := messageReceivedAt.UnixNano() / 1000000

				c.opts.OnTimeReceived(c, t0, t1, t2, t3)
			}
		case "stream":
			if c.opts.OnStreamCommand != nil {
				c.opts.OnStreamCommand(c, string(message.Payload()) == "on")
			}
		case "probespeed":
			if c.opts.OnProbeSpeedSet != nil {
				newspeed, err := strconv.ParseInt(string(message.Payload()), 10, 64)
				if err != nil {
					c.opts.Logger.WithError(err).Error("Error in PROBE SPEED payload")
				} else {
					c.opts.OnProbeSpeedSet(c, newspeed)
				}
			}
		default:
			// Unknown command
		}
	})

	return c.SendAlive()
}
