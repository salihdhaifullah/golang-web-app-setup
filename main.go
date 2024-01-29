package main

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"strings"
	"test/build"
	"text/template"
)

func gzipHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			w.Header().Set("Content-Encoding", "gzip")

			gzw := gzip.NewWriter(w)
			defer gzw.Close()

			w = gzipResponseWriter{gzw, w}
		}

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.gw.Write(b)
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func index_view(w http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("./pages/index.html")

	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title string
		Items []string
	}{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	build.Build()
	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", fs)

	router.HandleFunc("/", index_view)
	router.HandleFunc("/hello", hello)
	router.HandleFunc("/headers", headers)

	http.Handle("/", gzipHandler(router))

	log.Printf("server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
