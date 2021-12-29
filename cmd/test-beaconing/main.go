//nolint: forbidigo
package main

import (
	"context"
	"fmt"
	"git.sapienzaapps.it/SeismoCloud/seismocloud-client-go/localdiscovery"
)

func main() {
	fmt.Println("Waiting for discovery to reply")
	err := localdiscovery.Beacon(context.TODO(), "112233aabbcc", "esp8266", "1")
	if err != nil {
		panic(err)
	}
}
