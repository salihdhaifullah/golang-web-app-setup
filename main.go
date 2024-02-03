package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/a-h/templ"

	"github.com/salihdhaifullah/golang-web-app-setup/helpers"
	"github.com/salihdhaifullah/golang-web-app-setup/helpers/initializers"
	"github.com/salihdhaifullah/golang-web-app-setup/helpers/middleware"
	"github.com/salihdhaifullah/golang-web-app-setup/views"
)

func main() {
	initializers.GetENV()
	cmd := exec.Command("templ", "generate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	wg := sync.WaitGroup{};
	wg.Add(2)
	go helpers.WaitFor(initializers.MongoDB, &wg)
	if os.Getenv("ENV") == "PROD" {
		go helpers.WaitFor(helpers.Build, &wg)
	} else {
		go helpers.WaitFor(func () {
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}, &wg)
	}

	wg.Wait()

	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	mux.Handle("/", templ.Handler(views.Hello("salih dhaifullah")))

	http.Handle("/", middleware.Gzip(mux))

	initializers.Listen(mux)
}
