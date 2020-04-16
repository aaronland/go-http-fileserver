package main

import (
	"context"
	"flag"
	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-fileserver"
	"log"
	"net/http"
	"strings"
)

func main() {

	server_uri := flag.String("server-uri", "http://localhost:8080", "...")
	root := flag.String("root", "", "...")

	enable_cors := flag.Bool("enable-cors", false, "...")
	enable_gzip := flag.Bool("enable-gzip", false, "...")

	cors_origins := flag.String("cors-origins", "*", "...")

	flag.Parse()

	ctx := context.Background()

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Unable to create server (%s), %v", *server_uri, err)
	}

	if *root == "" {
		log.Fatalf("Missing -root")
	}

	fs_opts := &fileserver.FileServerOptions{
		Root:        *root,
		EnableCORS:  *enable_cors,
		CORSOrigins: strings.Split(*cors_origins, ","),
		EnableGzip:  *enable_gzip,
	}

	fs_handler, err := fileserver.NewFileServerHandler(fs_opts)

	if err != nil {
		log.Fatalf("Unable to create fileserver handler, %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", fs_handler)

	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
