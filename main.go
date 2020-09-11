package main

import (
	"crypto/tls"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

const port string = ":8080"

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
        Addr:         port,
        Handler:      serverHandler(),
        // For some reason disables http 2 "No clue how or why, must research.".
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
    }

    if proto == "http" {
        log.Printf("Server listening in http://localhost" + port)
        log.Fatal(server.ListenAndServe())
    } else {
        log.Printf("Server listening in https://localhost" + port)
        log.Fatal(server.ListenAndServeTLS(pemPath, keyPath))
    }

}

// Check if we use a tcp stream or just http
//For example, the CONNECT method can be used to access websites that use SSL (HTTPS).
//The client asks an HTTP Proxy server to tunnel the TCP connection 
//to the desired destination.
//The server then proceeds to make the connection on behalf of the client. 
//Once the connection has been established by the server, 
//the Proxy server continues to proxy the TCP stream to and from the client.
func serverHandler() http.HandlerFunc {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodConnect {
            handleStream(w, r)
        } else {
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
    // Try to connect to destination.
    log.Printf("Got tcp Stream request for %v\n", r.Host)
    destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
    if err != nil {
        // Return HTTP error response.
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }

    // Write ok response header.
    w.WriteHeader(http.StatusOK)

    hijacker, ok := w.(http.Hijacker)
    if !ok {
        http.Error(w, "Hijacking not supported.", http.StatusInternalServerError)
    }

    clientConn, _, err := hijacker.Hijack()
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }

    go transfer(clientConn, destConn)
    go transfer(destConn, clientConn)
}

func copyHeader(dest, src http.Header){
    for k, vv := range src {
        for _, v := range vv {
            dest.Add(k, v)
        }
    }
}

func handleHttp(w http.ResponseWriter, r *http.Request) {
    resp, err := http.DefaultTransport.RoundTrip(r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusServiceUnavailable)
        return
    }

    defer resp.Body.Close()
    w.WriteHeader(http.StatusOK)
    copyHeader(w.Header(), resp.Header)
    if resp.ContentLength != 0 {
        io.Copy(w, resp.Body)
    }
}

