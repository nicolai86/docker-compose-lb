package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	docker "github.com/fsouza/go-dockerclient"
)

func update(client *docker.Client, composeServiceName string) {
	options := docker.ListContainersOptions{
		Filters: map[string][]string{
			"label": {fmt.Sprintf("com.docker.compose.service=%v", composeServiceName)},
		},
	}

	cs, err := client.ListContainers(options)
	if err != nil {
		log.Printf("Failed to fetch container IDs: %v", err)
		return
	}

	mu.Lock()
	containers = cs
	mu.Unlock()
}

var (
	mu         *sync.Mutex = &sync.Mutex{}
	containers []docker.APIContainers
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":8080", "http address to listen on")
	flag.Parse()

	var composeServiceName = os.Getenv("DOCKER_COMPOSE_SERVICE_NAME")
	if composeServiceName == "" {
		log.Fatalf("Missing DOCKER_COMPOSE_SERVICE_NAME env variable")
	}

	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatalf("Failed to create docker client: %q", err)
	}

	update(client, composeServiceName)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			update(client, composeServiceName)
		}
	}()

	director := func(req *http.Request) {}
	proxy := &httputil.ReverseProxy{Director: director}
	http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		mu.Lock()
		target := containers[rand.Int()%len(containers)]
		mu.Unlock()

		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("%s:%d", target.Networks.Networks["bridge"].IPAddress, target.Ports[0].PrivatePort)
		w.Header().Set("X-Internal-Service", req.URL.Host)

		proxy.ServeHTTP(w, req)
	}))
}
