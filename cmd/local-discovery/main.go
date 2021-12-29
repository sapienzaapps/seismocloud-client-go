//nolint:forbidigo
package main

import (
	"context"
	"fmt"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-client-go/localdiscovery"
)

func main() {
	fmt.Println("Starting local discovery...")
	var sensors = make(chan localdiscovery.Sensor, 1)
	go func() {
		for c := range sensors {
			fmt.Println(c)
		}
	}()
	_, err := localdiscovery.Scan(context.TODO(), sensors)
	if err != nil {
		panic(err)
	}
}
