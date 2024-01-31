package main

import (
	"compress/gzip"
	"log"
	"net/http"
	"os"
	"strings"
	"test/api"
	"test/build"
	"text/template"

	"github.com/nickalie/go-webpbin"
)

type gzipResponseWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.gw.Write(b)
}

func handler(w http.ResponseWriter, r *http.Request) {
	img := api.CreateImage()
	f, e := os.Create("./test.webp")
	if e != nil {
		log.Fatal(e)
	}

	webpbin.Encode(f, img)

	w.Header().Set("Content-Type", "image/webp")

	err := webpbin.Encode(w, img)

	if err != nil {
		log.Fatal(err)
	}
}

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
	router.HandleFunc("/img", handler)

	http.Handle("/", gzipHandler(router))

	log.Printf("server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
