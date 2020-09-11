# Golang implemetation of a http proxy.

## To run 
    - `go get github.com/dript0hard/goproxy`
    - `./setup.sh` Generate the cert and key for https.
    - `go run main.go` Port 8080
    - `curl -Lv https://www.instagram.com/deni_myftiu --proxy https://localhost:8080 --proxy-cacert server.pem` try ro curl instagram for a tls connection.
