package build

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/salihdhaifullah/golang-web-app-setup/helpers/image_processor"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

func initMinify() *minify.M {
	var minifyer = minify.New()

	minifyer.AddFunc("text/css", css.Minify)
	minifyer.AddFunc("image/svg+xml", svg.Minify)
	minifyer.AddFunc("text/javascript", js.Minify)

	return minifyer
}

var minifyer = initMinify()

var contentTypes = map[string]string{
	".js":   "text/javascript",
	".css":  "text/css",
	".svg":  "image/svg+xml",
}

func visitFile(fp string, fi os.DirEntry, err error) error {
	if err != nil {
		log.Println(err)
		return nil
	}

	if fi.IsDir() {
		return nil
	}

	buildPath := filepath.Join("./dist", fp)
	dir := filepath.Dir(buildPath)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	w, e := os.Create(buildPath)

	if e != nil {
		log.Printf("from create func %s\n", buildPath)
		log.Fatal(e)
	}

	r, e := os.Open(fp)

	if e != nil {
		log.Printf("from open func %s\n", fp)
		log.Fatal(e)
	}

	defer w.Close()
	defer r.Close()

	ext := filepath.Ext(fp)
	FileType, ok := contentTypes[ext]

	if !ok {
		log.Printf("coping file from %s to %s", fp, buildPath)
		_, err = io.Copy(w, r)
		if err != nil {
			return err
		}
	} else if err := minifyer.Minify(FileType, w, r); err != nil {
		log.Fatal(err)
	}

	return nil
}

func listFiles(rootDir string) {
	err := filepath.WalkDir(rootDir, visitFile)
	if err != nil {
		log.Printf("Error walking the path %v: %v\n", rootDir, err)
	}
}

func Build() {
	e := image_processor.GenerateResizedImages("./static/test.png", "./")
	if e != nil {
		log.Fatal(e)
	}

	listFiles("./static")
}
