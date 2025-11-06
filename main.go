package main

import (
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
)


type apiConfig struct{
	fileserverHits atomic.Int32
}

func main(){
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz", handlerHealthz)
	mux.HandleFunc("GET /metrics", apiCfg.handlerHits)
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

    log.Println("Starting server on :8080")
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

func handlerHealthz(resWriter http.ResponseWriter,req *http.Request){
	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)

	_, err := resWriter.Write([]byte("OK"))
	if err != nil{
		log.Println(err)
	}
}

func (apiCfg *apiConfig) handlerHits(resWriter http.ResponseWriter,req *http.Request){
	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)


    hitsStr := "Hits: " + strconv.Itoa(int(apiCfg.fileserverHits.Load()))
    _, err := resWriter.Write([]byte(hitsStr))
    if err != nil {
        log.Println(err)
    }
}

func (apiCfg *apiConfig) handlerReset(resWriter http.ResponseWriter,req *http.Request){
	resWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	resWriter.WriteHeader(200)

	apiCfg.fileserverHits.Store(0)
    hitsStr := "Hits: " + strconv.Itoa(int(apiCfg.fileserverHits.Load()))
    _, err := resWriter.Write([]byte(hitsStr))
    if err != nil {
        log.Println(err)
    }
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cfg.fileserverHits.Add(1)
        next.ServeHTTP(w, r)
    })
}