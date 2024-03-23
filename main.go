package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"text/template"
)

type ServerData struct {
	URL               string
	RequestsProcessed int
}

type StatsData struct {
	TotalRequestCount int
	ServersData       []ServerData
}

type Server struct {
	URL               string
	RequestsProcessed int
	Mutex             sync.Mutex
}

func (s *Server) parseURL() (*url.URL, error) {
	url, err := url.Parse(s.URL)
	if err != nil {
		return nil, err
	}
	return url, nil
}

type LoadBalancer struct {
	Servers      []*Server
	RequestCount int
	Mutex        sync.Mutex
}

func (lb *LoadBalancer) RoundRobin() *Server {
	nextServerIndex := lb.RequestCount % len(lb.Servers)
	return lb.Servers[nextServerIndex]
}

var servers = []*Server{
	{URL: "http://localhost:9770"},
	{URL: "http://localhost:9771"},
	{URL: "http://localhost:9772"},
	{URL: "http://localhost:9773"},
	{URL: "http://localhost:9774"},
}

func main() {
	loadBalancer := &LoadBalancer{Servers: servers}

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.Mutex.Lock()
		loadBalancer.RequestCount++
		loadBalancer.Mutex.Unlock()

		nextServer := loadBalancer.RoundRobin()

		url, err := nextServer.parseURL()
		if err != nil {
			fmt.Println("could not parse server url", err)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.ServeHTTP(w, r)

		nextServer.Mutex.Lock()
		nextServer.RequestsProcessed++
		nextServer.Mutex.Unlock()
	})

	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		data := StatsData{
			TotalRequestCount: loadBalancer.RequestCount,
			ServersData:       []ServerData{},
		}

		for _, s := range loadBalancer.Servers {
			data.ServersData = append(data.ServersData, ServerData{
				URL:               s.URL,
				RequestsProcessed: s.RequestsProcessed,
			})
		}

		isRequestFromHTMX := r.Header.Get("HX-Request") == "true"

		// just update the content
		if isRequestFromHTMX {
			tmpl := template.Must(template.New("content").ParseFiles("./views/index.html"))
			tmpl.ExecuteTemplate(w, "content", data)
			return
		}

		tmpl := template.Must(template.New("index").ParseFiles("./views/index.html"))
		tmpl.ExecuteTemplate(w, "index", data)
	})

	fmt.Println("load balancer running at 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
