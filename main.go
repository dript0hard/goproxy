package main

import (
	"flag"
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn, dest *string) {
    defer conn.Close()
    connTo, err := net.Dial("tcp", *dest + ":443")
    if err != nil {
        log.Fatal(err)
        return
    }
    defer connTo.Close()
    go io.Copy(connTo, conn)
    io.Copy(conn, connTo)
}

func main(){
    var siteToProxy string
    flag.StringVar(&siteToProxy, "proxyTo", "localhost", "What do we want to proxy.")
    flag.Parse()

    l, err := net.Listen("tcp", "0.0.0.0:9001")
    if err != nil {
        log.Fatal(err)
    }

    for {
        conn, err := l.Accept()
        log.Print("Conn accepted.")
        if err != nil {
            log.Fatal(err)
        }
        go handleConnection(conn, &siteToProxy)
    }
}
