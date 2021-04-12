package fileserver

import (
	"errors"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/rs/cors"
	"net/http"
	"os"
	"path/filepath"
)

type FileServerOptions struct {
	Root        string
	EnableCORS  bool
	CORSOrigins []string
	EnableGzip  bool
}

type ContentTypeOptions struct {
	Matches map[string]string
}

func NewContentTypeHandler(opts *ContentTypeOptions, next http.Handler) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		uri := req.RequestURI
		ext := filepath.Ext(uri)

		for re_ext, content_type := range opts.Matches {

			if ext == re_ext {
				rsp.Header().Set("Content-Type", content_type)
				break
			}
		}

		next.ServeHTTP(rsp, req)
	}

	return http.HandlerFunc(fn), nil
}

func NewFileServerHandler(opts *FileServerOptions) (http.Handler, error) {

	abs_root, err := filepath.Abs(opts.Root)

	if err != nil {
		return nil, err
	}

	info, err := os.Stat(abs_root)

	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		msg := fmt.Sprintf("Root (%s) is not a directory", abs_root)
		return nil, errors.New(msg)
	}

	http_root := http.Dir(abs_root)
	fs_handler := http.FileServer(http_root)

	if opts.EnableGzip {
		fs_handler = gziphandler.GzipHandler(fs_handler)
	}

	if opts.EnableCORS {

		c := cors.New(cors.Options{
			AllowedOrigins: opts.CORSOrigins,
		})

		fs_handler = c.Handler(fs_handler)
	}

	return fs_handler, nil
}
