package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

func main(){
    var pemPath string
    flag.StringVar(&pemPath, "pem", "server.pem", "Path to pem file.")

    var keyPath string
    flag.StringVar(&keyPath, "key", "server.key", "Path to key for server.pem.")

    var proto string
    flag.StringVar(&proto, "proto", "https", "Protocoll for the proxy either http or https.")

    flag.Parse()

    if proto != "https" && proto != "http" {
        log.Fatal("Unrecognised protocoll. Must be either http or https.")
    }

    server := &http.Server{
        Addr:         ":8080",
        Handler:      serverHandler(),
        // For some reason disables http 2 "No clue how or why, must research.".
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
    }

    if proto == "http" {
        log.Printf("Server listening in http://localhost:8888")
        log.Fatal(server.ListenAndServe())
    } else {
        log.Printf("Server listening in https://localhost:8888")
        log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
    }

}

func serverHandler() http.HandlerFunc {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if we use a tcp stream or just http
        if r.Method == http.MethodConnect {
            handleStream(w, r)
        } else {
            handleHttp(w, r)
        }
    })

    return handler
}

func handleHttp(w http.ResponseWriter, r *http.Request) {

}

func handleStream(w http.ResponseWriter, r *http.Request) {

}
