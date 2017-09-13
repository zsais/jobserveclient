package main

import (
	"fmt"

	"github.com/zsais/jobserveclient/jsclient"
)

func main() {
	l, err := jsclient.NewListener()
	if err != nil {
		fmt.Printf("error creating NewListener(): %v\n", err)
	}
	jsclient.ProcessConnections(l)
}
