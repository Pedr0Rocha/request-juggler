package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	fmt.Printf("API running at %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", "8000"), nil)
	if err != nil {
		panic(fmt.Sprintf("server %s crashed", port))
	}
}
