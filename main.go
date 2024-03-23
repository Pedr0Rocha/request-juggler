package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Server struct {
	URL               *url.URL
	RequestsProcessed int
	Mutex             sync.Mutex
}

type LoadBalancer struct {
	Servers      []*Server
	RequestCount int
	Mutex        sync.Mutex
}

func NewLoadBalancer() *LoadBalancer {
	serverUrls := []string{
		"http://localhost:9770",
		"http://localhost:9771",
		"http://localhost:9772",
		"http://localhost:9773",
		"http://localhost:9774",
	}

	var servers []*Server
	for _, s := range serverUrls {
		url, _ := url.Parse(s)
		servers = append(servers, &Server{
			URL: url,
		})
	}

	return &LoadBalancer{
		Servers: servers,
	}
}

func (lb *LoadBalancer) RoundRobin() *Server {
	nextServerIndex := lb.RequestCount % len(lb.Servers)
	return lb.Servers[nextServerIndex]
}

var loadBalancer = NewLoadBalancer()

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.Mutex.Lock()
		nextServer := loadBalancer.RoundRobin()
		loadBalancer.RequestCount++
		loadBalancer.Mutex.Unlock()

		proxy := httputil.NewSingleHostReverseProxy(nextServer.URL)
		proxy.ServeHTTP(w, r)

		nextServer.Mutex.Lock()
		nextServer.RequestsProcessed++
		nextServer.Mutex.Unlock()
	})

	http.HandleFunc("/stats", StatsHandler)
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	fmt.Println("load balancer running at 8080")
	http.ListenAndServe(":8080", nil)
}
