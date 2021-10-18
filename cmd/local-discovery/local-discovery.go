package main

import (
	"context"
	"fmt"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-client-go/localdiscovery"
)

func main() {
	var sensors = make(chan localdiscovery.Sensor, 1)
	go func() {
		for c := range sensors {
			fmt.Println(c)
		}
	}()
	_, err := localdiscovery.Discovery(context.TODO(), sensors)
	if err != nil {
		panic(err)
	}
}
