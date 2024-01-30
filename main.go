package main

import (
	"compress/gzip"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"net/http"
	"strings"
	"test/build"
	"text/template"
)

func createImage() *image.RGBA {
	width := 800
	height := 600
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	blue := color.RGBA{0, 0, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Point{}, draw.Src)

	red := color.RGBA{255, 0, 0, 255}
	centerX, centerY := width/2, height/2
	radius := 100

	for x := -radius; x < radius; x++ {
		for y := -radius; y < radius; y++ {
			if x*x+y*y < radius*radius {
				img.Set(centerX+x, centerY+y, red)
			}
		}
	}

	return img
}

func handler(w http.ResponseWriter, r *http.Request) {
	img := createImage()

	w.Header().Set("Content-Type", "image/jpeg")

	err := jpeg.Encode(w, img, nil)

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
	router.HandleFunc("/img", handler)

	http.Handle("/", gzipHandler(router))

	log.Printf("server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
