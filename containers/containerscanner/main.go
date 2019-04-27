package main

import (
	"os/signal"
	"os"
	"sync"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	whale "github.com/docker/docker/client"
	"strconv"
	"time"
	"regexp"
	"strings"
)
// ------------
// plugins
// ------------

// Endpoint represent the service endpoint
type Endpoint struct {
	clustername    string
	port           uint16
	host           string
	frontProxyPath string
	version        string
}

// EndpointUpdateRequest represent the update request
type EndpointUpdateRequest struct {
	pluginName string
	timestamp  time.Time
	endpoints  []Endpoint
}

// Plugin is the extension point for configuration data sources
type Plugin interface {
	getName() string
	run(ctx context.Context) chan *EndpointUpdateRequest
}

func RunAllPlugins(ctx context.Context, plugins []Plugin) chan *EndpointUpdateRequest {
	channel := make (chan *EndpointUpdateRequest)
	go func(){
		defer close(channel)
		var waitgroup sync.WaitGroup
		for _, plugin := range plugins {
			waitgroup.Add(1)
			go func(p *Plugin, wg *sync.WaitGroup){
				defer wg.Done() 
				pluginChannel := (*p).run(ctx)
				for {
				  request, more := <-pluginChannel
				  if !more {
					fmt.Printf("[%s] fan in loop terminating, channel closed\n", (*p).getName())  
					return
				  }
				  channel<-request
				} 
			}(&plugin, &waitgroup)		
		}
		waitgroup.Wait()
		fmt.Println("[Plugins] All plugins terminated")  
	}()
	return channel	
}

// ---------------
// docker plugin
// ---------------

// Docker provides configuration source from docker
type Docker struct {
}

const (
	docker      = "Docker"
	commitIDKey = "COMMIT_ID"
	versionKey  = "VERSION"
)

var (
	portGroupExpr          = "(?P<port>\\d+)"
	urlPrefixExpr          = fmt.Sprintf("CLUSTER_%s_URLPREFIX", portGroupExpr)
	urlPrefixPattern       = regexp.MustCompile(urlPrefixExpr)
	serviceNameExpr        = fmt.Sprintf("CLUSTER_%s_NAME", portGroupExpr)
	serviceNamePattern     = regexp.MustCompile(serviceNameExpr)
	serviceCategoryExpr    = fmt.Sprintf("CLUSTER_%s_CATEGORY", portGroupExpr)
	serviceCategoryPattern = regexp.MustCompile(serviceCategoryExpr)
)

var (
	serviceNameSubmatchGroupLookup = ToMap(serviceNamePattern.SubexpNames())
	portIndex                      = serviceNameSubmatchGroupLookup["port"]
)

type service struct {
	name      string
	category  string
	urlPrefix string
	version   string
	port      uint16
}

type discoverableContainer struct {
	services []service
	types.Container
}

func  getServicePorts(container types.Container) []uint16 {
	ports := []uint16{}
	uniquePortsContainer := make(map[uint16] string)
	for key := range container.Labels {
		if serviceNamePattern.MatchString(key) {
			submatches := serviceNamePattern.FindStringSubmatch(key)
			port:= uint16(MustAtoi(submatches[portIndex]))
			if _, alreadyExists := uniquePortsContainer[port]; !alreadyExists {
				ports = append(ports, port)
				uniquePortsContainer[port] = "exist"
			}
		}
	}
	return ports
}

func mapContainerToDiscoverableContainer(container types.Container, servicePorts []uint16) *discoverableContainer {
	labels := container.Labels
	discoveredContainer := &discoverableContainer{
		Container: container,
	}
	services := []service{}
	for _, port := range servicePorts {
		serviceNameLabelKey := strings.Replace(serviceNameExpr, portGroupExpr, strconv.Itoa(int(port)), 1)
		serviceCategoryLabelKey := strings.Replace(serviceCategoryExpr, portGroupExpr, strconv.Itoa(int(port)), 1)
		urlPrefixLabelKey := strings.Replace(urlPrefixExpr, portGroupExpr, strconv.Itoa(int(port)), 1)

		service := service{
			name:      labels[serviceNameLabelKey],
			category:  labels[serviceCategoryLabelKey],
			urlPrefix: labels[urlPrefixLabelKey],
			version:   fmt.Sprintf("v%s-%s", labels[versionKey], labels[commitIDKey]),
			port:      port,
		}
		services = append(services, service)
	}
	discoveredContainer.services = services
	return discoveredContainer
}

func getDiscoverableContainers(containers []types.Container) *[]discoverableContainer {
	discoveredContainers := []discoverableContainer{}
	for _, container := range containers {
		servicePorts := getServicePorts(container)
		if len(servicePorts) > 0 {
			discoveredContainers = append(discoveredContainers,
				*mapContainerToDiscoverableContainer(container, servicePorts))
		}
	}
	return &discoveredContainers
}

type enPorts []types.Port
func (ports enPorts) wherePorts(predicate func(types.Port) bool) []types.Port {
	matchingPorts := []types.Port{}
	for _, port := range ports {
		if predicate(port) {
			matchingPorts = append(matchingPorts, port)
		}
	}
	return matchingPorts
}

func (ports enPorts) getMappedAddress(portNumber uint16) (mappedHost string, mappedPortNumber uint16) {
	mappedPorts := ports.
		wherePorts(func(p types.Port) bool {
			return p.PrivatePort == portNumber
		})
	if len(mappedPorts) > 0 {
		mappedHost = "host.docker.internal"
		mappedPortNumber = mappedPorts[0].PublicPort
	}
	return
}

func (container *discoverableContainer) mapToEndpointUpdateRequest() *EndpointUpdateRequest {
	endpoints := []Endpoint{}
	for _, service := range container.services {
		host, portNumber := enPorts(container.Ports).
			getMappedAddress(service.port)

		frontProxyPath := fmt.Sprintf("/%s/%s", service.category,
			service.name)
		if service.urlPrefix != "" {
			frontProxyPath = fmt.Sprintf("/%s/%s", service.category,
				service.urlPrefix)
		}

		endpoint := Endpoint{
			clustername:    service.name,
			host:           host,
			port:           portNumber,
			frontProxyPath: frontProxyPath,
			version:        service.version,
		}
		endpoints = append(endpoints, endpoint)
	}

	request := &EndpointUpdateRequest{
		pluginName: docker,
		timestamp:       time.Now(),
		endpoints:  endpoints,
	}

	return request
}

func (plugin *Docker) getName() string{
	return docker
}

func (plugin *Docker) run(ctx context.Context) chan *EndpointUpdateRequest{
	cli, err := whale.NewEnvClient()
	if err != nil {
		panic(err)
	}
	channel := make(chan *EndpointUpdateRequest)
	go func (plugin *Docker){
		defer close(channel)
		for {

			containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
			if err != nil {
				panic(err)
			}

			discoveredContainers := getDiscoverableContainers(containers)

			for _, container := range *discoveredContainers {
				select {
				case channel <-container.mapToEndpointUpdateRequest():
					fmt.Printf("[%s] sending update request\n", plugin.getName())
				case <-ctx.Done():
					fmt.Printf("[%s] terminating scanner loop\n", plugin.getName())
					return		
				}
				
			}

			select {
			case <-time.After(2 * time.Second):
				fmt.Printf("[%s] going again\n", plugin.getName())
			case <-ctx.Done():
				fmt.Printf("[%s] terminating scanner loop\n", plugin.getName())
				return
			}
		}
	}(plugin)
	return channel
}

// -------------------
// utility
// ---------------------

// ToMap maps string array to map[string] int
func ToMap(texts []string) map[string]int {
	m := make(map[string]int)

	for index, text := range texts {
		m[text] = index
	}
	return m
}

// MustAtoi returns an integer, throws an exception if there is error
func MustAtoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

// -----------------
func main(){
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	defer func() {
		signal.Stop(c)
		cancel()
		fmt.Println("[main] the end")
	}()
	go func() {
		select {
		case <-c:
			cancel()
		case <-ctx.Done():
		}
	}()

	dockerPlugin := &Docker{
	}
	allPlugins := []Plugin{dockerPlugin}
	channel := RunAllPlugins(ctx, allPlugins)
	for {
		select {
		case updateRequest := <-channel:
			fmt.Printf("%v\n", updateRequest)

		case <-ctx.Done():
			return
		}
	}
}