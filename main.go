package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/a-h/templ"

	"github.com/salihdhaifullah/golang-web-app-setup/helpers"
	"github.com/salihdhaifullah/golang-web-app-setup/helpers/initializers"
	"github.com/salihdhaifullah/golang-web-app-setup/helpers/middleware"
	"github.com/salihdhaifullah/golang-web-app-setup/views"
)

func main() {
	initializers.GetENV()
	wg := sync.WaitGroup{};
	wg.Add(1)


	go helpers.WaitFor(initializers.MongoDB, &wg)
	if os.Getenv("ENV") == "PROD" {
		wg.Add(1)
		go helpers.WaitFor(helpers.Build, &wg)
	}

	wg.Wait()

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/", templ.Handler(views.Hello("salih dhaifullah")))

	http.Handle("/", middleware.Gzip(mux))
	initializers.Listen(mux)
}
