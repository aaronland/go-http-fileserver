package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-http-fileserver"
	"github.com/aaronland/go-http-server"
	"github.com/sfomuseum/go-flags/multi"
)

func main() {

	schemes := server.Schemes()
	str_schemes := strings.Join(schemes, ", ")

	server_desc := fmt.Sprintf("A valid aaronland/go-http-server URI. Registered schemes are: %s.", str_schemes)

	server_uri := flag.String("server-uri", "http://localhost:8080", server_desc)
	root := flag.String("root", "", "A valid path to serve files from.")
	prefix := flag.String("prefix", "", "A prefix to append to URL to serve requests from.")

	enable_cors := flag.Bool("enable-cors", false, "Enable CORS headers on responses.")
	enable_gzip := flag.Bool("enable-gzip", false, "Enable gzip-ed responses.")

	cors_origins := flag.String("cors-origins", "*", "A comma-separated of origins to allow CORS requests from.")

	var mimetypes multi.KeyValueString
	flag.Var(&mimetypes, "mimetype", "One or more key=value pairs mapping a file extension to a specific content (or mime) type to assign for that request")

	flag.Parse()

	ctx := context.Background()

	s, err := server.NewServer(ctx, *server_uri)

	if err != nil {
		log.Fatalf("Unable to create server (%s), %v", *server_uri, err)
	}

	if *root == "" {

		cwd, err := os.Getwd()

		if err != nil {
			log.Fatalf("Failed to derive current working directory, %v", err)
		}

		*root = cwd
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

	root_handler := fs_handler

	//

	if len(mimetypes) > 0 {

		matches := make(map[string]string)

		for _, kv := range mimetypes {

			matches[kv.Key()] = kv.Value().(string)
		}

		ct_opts := &fileserver.ContentTypeOptions{
			Matches: matches,
		}

		ct_handler, err := fileserver.NewContentTypeHandler(ct_opts, fs_handler)

		if err != nil {
			log.Fatalf("Failed to create new content type handler, %v", err)
		}

		root_handler = ct_handler
	}

	//

	mux := http.NewServeMux()
	mux.Handle(uri, root_handler)

	log.Printf("Serving %s and listening for requests on %s", *root, s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		log.Fatalf("Failed to start server, %v", err)
	}
}
