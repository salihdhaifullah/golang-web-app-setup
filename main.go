package main

import (
	"compress/gzip"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/nickalie/go-webpbin"

	"github.com/salihdhaifullah/golang-web-app-setup/build"
	"github.com/salihdhaifullah/golang-web-app-setup/helpers/image_processor"
	"github.com/salihdhaifullah/golang-web-app-setup/helpers/initializers"
	"github.com/salihdhaifullah/golang-web-app-setup/views"
)

type gzipResponseWriter struct {
	gw *gzip.Writer
	http.ResponseWriter
}

func (grw gzipResponseWriter) Write(b []byte) (int, error) {
	return grw.gw.Write(b)
}

func handler(w http.ResponseWriter, r *http.Request) {
	img := image_processor.CreateImage()
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

func main() {
	initializers.GetENV()

	if os.Getenv("env") != "DEV" {
		build.Build()
	}

	router := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	router.Handle("/", templ.Handler(views.Hello("salih dhaifullah")))
	router.HandleFunc("/img", handler)

	http.Handle("/", gzipHandler(router))

	log.Printf("server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
