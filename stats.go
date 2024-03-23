package main

import (
	"net/http"
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

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	data := StatsData{
		TotalRequestCount: loadBalancer.RequestCount,
		ServersData:       []ServerData{},
	}

	for _, s := range loadBalancer.Servers {
		data.ServersData = append(data.ServersData, ServerData{
			URL:               s.URL.String(),
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
}
