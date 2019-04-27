package main

import (
	"fmt"
	"context"
	"github.com/docker/docker/api/types"
	whale "github.com/docker/docker/client"
)
type enPorts []types.Port;
func (ports enPorts) wherePorts(predicate func(types.Port) bool ) []types.Port{
	matchingPorts := []types.Port{}
	for _ ,port := range ports {
		if predicate(port) {
			matchingPorts = append(matchingPorts, port)
		}
	}
	return matchingPorts
}

func (ports enPorts) getMappedAddress (portNumber uint16) (mappedIp string, mappedPortNumber uint16) {
	mappedPorts := ports.
					wherePorts(func(p types.Port)bool {
						return p.PrivatePort == portNumber
					})
	if len(mappedPorts)> 0 {
		mappedIp = mappedPorts[0].IP
		mappedPortNumber = mappedPorts[0].PublicPort
	}
	return
}

func main() {
	cli, err := whale.NewClientWithOpts(whale.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		mappedPorts := enPorts(container.Ports).
			wherePorts(func(port types.Port) bool{
				return port.PrivatePort == 80
			})		
		fmt.Println(container.ID)
		if len(mappedPorts) > 0 {
			fmt.Printf( "80 mapped to %d\n" , mappedPorts[0].PublicPort)
		}
	}

	
}