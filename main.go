package main

import (
	"crypto/tls"
	"flag"
	"io"
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
        //For example, the CONNECT method can be used to access websites that use SSL (HTTPS).
        //The client asks an HTTP Proxy server to tunnel the TCP connection 
        //to the desired destination.
        //The server then proceeds to make the connection on behalf of the client. 
        //Once the connection has been established by the server, 
        //the Proxy server continues to proxy the TCP stream to and from the client.
        if r.Method == http.MethodConnect {
            PrintRequest(r)
            handleStream(w, r)
        } else {
            PrintRequest(r)
            handleHttp(w, r)
        }
    })

    return handler
}

func transfer(dest io.WriteCloser, src io.ReadCloser) {
    defer dest.Close()
    defer src.Close()
    io.Copy(dest, src)
}

func handleStream(w http.ResponseWriter, r *http.Request) {

}

func handleHttp(w http.ResponseWriter, r *http.Request) {

}

