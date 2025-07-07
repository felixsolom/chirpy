package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const rootDir = "."

	mux := http.NewServeMux()
	mux.Handle("/app", http.FileServer(http.Dir(rootDir)))
	mux.Handle("/assets", http.FileServer(http.Dir(rootDir+"/assets")))

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %v", port)
	log.Fatal(srv.ListenAndServe())
}
