package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareNumOfRequests() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits := cfg.fileserverHits.Load()
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)
	})
}

func (cfg *apiConfig) middlewareResetReqCounter() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0"))
	})
}

func handlerHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Error string `json:"error"`
		Valid bool   `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding parameters %s", err)
		respBody := returnVals{
			Error: "Something went wrong",
			Valid: false,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		respBody := returnVals{
			Error: "Chirp is too long",
			Valid: false,
		}
		dat, err := json.Marshal(respBody)
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	respBody := returnVals{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func main() {
	const port = "8080"
	const rootDir = "."
	cfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(rootDir)))))
	mux.Handle("GET /api/healthz", http.StripPrefix("/api/", http.HandlerFunc(handlerHealth)))
	mux.Handle("POST /api/validate_chirp", http.HandlerFunc(validateChirp))
	mux.Handle("GET /admin/metrics", cfg.middlewareNumOfRequests())
	mux.Handle("POST /admin/reset", cfg.middlewareResetReqCounter())

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving on port: %v", port)
	log.Fatal(srv.ListenAndServe())
}
