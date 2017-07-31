package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.String("port", "8000", "http listen port")
	flag.Parse()

	fmt.Println("Serving HTTP on 0.0.0.0 port", *port)

	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":"+*port, nil)
}
