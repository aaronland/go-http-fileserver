# go-http-fileserver

There are many Go HTTP file server tools. This one is ours.

## Important

Documentation to follow.

## Tools

### fileserver

```
$> ./bin/fileserver -h
Usage of ./bin/fileserver:
  -cors-origins string
    	A comma-separated of origins to allow CORS requests from. (default "*")
  -enable-cors
    	Enable CORS headers on responses.
  -enable-gzip
    	Enable gzip-ed responses.
  -mimetype value
    	One or more key=value pairs mapping a file extension to a specific content (or mime) type to assign for that request
  -prefix string
    	A prefix to append to URL to serve requests from.
  -root string
    	A valid path to serve files from.
  -server-uri string
    	A valid aaronland/go-http-server URI. Registered schemes are: HTTP,HTTPS,LAMBDA,MKCERT,TLS. (default "http://localhost:8080")
```

## See also

* https://github.com/aaronland/go-http-server