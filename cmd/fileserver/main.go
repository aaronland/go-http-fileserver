package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-http-fileserver"
	"github.com/aaronland/go-http-server"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {

	schemes := server.SchemesAsString()

	server_desc := fmt.Sprintf("A valid aaronland/go-http-server URI. Registered schemes are: %s.", schemes)

	server_uri := flag.String("server-uri", "http://localhost:8080", server_desc)
	root := flag.String("root", "", "A valid path to serve files from.")
	prefix := flag.String("prefix", "", "A prefix to append to URL to serve requests from.")

	enable_cors := flag.Bool("enable-cors", false, "Enable CORS headers on responses.")
	enable_gzip := flag.Bool("enable-gzip", false, "Enable gzip-ed responses.")

	cors_origins := flag.String("cors-origins", "*", "A comma-separated of origins to allow CORS requests from.")

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

	uri := "/"

	if *prefix != "" {

		if !strings.HasPrefix(*prefix, "/") {
			log.Fatalf("Invalid prefix. Missing leading /")
		}

		fs_handler = http.StripPrefix(*prefix, fs_handler)
		uri = filepath.Join(uri, *prefix)
	}

	if !strings.HasSuffix(uri, "/") {
		uri = fmt.Sprintf("%s/", uri)
	}

	mux := http.NewServeMux()
	mux.Handle(uri, fs_handler)

	log.Printf("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
