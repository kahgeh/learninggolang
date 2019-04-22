package main

import (
	"fmt"
	"net/http"
	"time"
)

type MiddlewareMessage struct {
	source string
	values map[string]interface{}
}

type Middleware struct {
	name           string
	prevMiddleware *Middleware
	nextMiddleware *Middleware
	outputChannel  chan MiddlewareMessage
	handler        func(*http.ResponseWriter, *http.Request, MiddlewareMessage) map[string]interface{}
}

type Route struct {
	routeTemplate string
	handler       func(*http.ResponseWriter, *http.Request, MiddlewareMessage)
}

type RouteRegister struct {
	routes []Route
}

func (register *RouteRegister) Register(route Route) *RouteRegister {
	(*register).routes = append(register.routes, route)
	return register
}

func merge(maps ...map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})

	for _, mp := range maps {
		if mp != nil {
			for key, value := range mp {
				newMap[key] = value
			}
		}
	}

	return newMap
}

func (middleware *Middleware) GetDoneChannel() chan MiddlewareMessage {
	curMiddleware := middleware
	for {
		if curMiddleware == nil {
			return nil
		}

		if curMiddleware.nextMiddleware == nil {
			return curMiddleware.outputChannel
		}

		curMiddleware = curMiddleware.nextMiddleware
	}
}

func (middleware *Middleware) AddMiddleware(name string, middlewareHandler func(*http.ResponseWriter, *http.Request, MiddlewareMessage) map[string]interface{}) *Middleware {
	newMiddleware := &Middleware{
		name:           name,
		handler:        middlewareHandler,
		outputChannel:  make(chan MiddlewareMessage),
		prevMiddleware: middleware,
	}
	middleware.nextMiddleware = newMiddleware
	fmt.Printf("%s address = %p\n", newMiddleware.name, &newMiddleware)
	return newMiddleware
}

func (middleware *Middleware) Run(w *http.ResponseWriter, request *http.Request) {
	fmt.Printf("Running %s (%p) ...\n", middleware.name, middleware)
	input := MiddlewareMessage{}
	if middleware.prevMiddleware != nil {
		fmt.Println("reading previous middleware output")
		input = <-middleware.prevMiddleware.outputChannel
	}
	fmt.Printf("running %s handler\n", middleware.name)
	result := middleware.handler(w, request, input)
	output := MiddlewareMessage{
		values: merge(input.values, result),
	}

	if middleware.nextMiddleware != nil {
		fmt.Println("kickoff next middleware")
		go middleware.nextMiddleware.Run(w, request)
		middleware.outputChannel <- output
		fmt.Printf("Completed %s\n", middleware.name)
		return
	}

	if middleware.outputChannel != nil {
		middleware.outputChannel <- output
	}

	fmt.Printf("Completed %s\n", middleware.name)
}

func authenticationMiddlewareHandler(w *http.ResponseWriter, request *http.Request, input MiddlewareMessage) map[string]interface{} {
	outputValues := make(map[string]interface{})
	authToken := request.Header.Get("Authorization")
	time.Sleep(1 * time.Second)
	if authToken != "" {
		outputValues["userid"] = "testuserid"
		return outputValues
	}
	return outputValues
}

func loggingMiddleware(w *http.ResponseWriter, request *http.Request, input MiddlewareMessage) map[string]interface{} {
	fmt.Printf("%s Received request on /%s\n", time.Now().Format(time.RFC3339), request.URL.Path[1:])
	return nil
}

func makeResourceMiddleware() (routeRegister *RouteRegister, handler func(*http.ResponseWriter, *http.Request, MiddlewareMessage) map[string]interface{}) {
	routeRegister = new(RouteRegister)
	routeRegister.routes = []Route{}

	handler = func(responseWriter *http.ResponseWriter, request *http.Request, input MiddlewareMessage) map[string]interface{} {
		fmt.Println("resolving handler ...")
		routeRegister.routes[0].handler(responseWriter, request, input)
		return nil
	}
	return
}

func staticFileHandler(writer *http.ResponseWriter, request *http.Request, input MiddlewareMessage) {
	fmt.Fprintf(*writer, "<html><body>hello %s</body></html>", input.values["userid"])
}

func setupPipeline() *Middleware {
	firstMiddleware := &Middleware{
		name:          "logging",
		handler:       loggingMiddleware,
		outputChannel: make(chan MiddlewareMessage),
	}
	fmt.Printf("%s address = %p\n", firstMiddleware.name, &firstMiddleware)
	routeRegister, resourceMiddlewareHandler := makeResourceMiddleware()
	routeRegister.Register(Route{routeTemplate: "*.html", handler: staticFileHandler})
	firstMiddleware.
		AddMiddleware("authentication", authenticationMiddlewareHandler).
		AddMiddleware("resource", resourceMiddlewareHandler)

	return firstMiddleware
}

func main() {
	head := setupPipeline()
	http.HandleFunc("/", func(w http.ResponseWriter, request *http.Request) {
		go head.Run(&w, request)
		<-head.GetDoneChannel()
	})
	http.ListenAndServe(":8080", nil)
}
