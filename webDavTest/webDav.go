package main

import (
	"net/http"
	"os"

	// "flag"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/net/webdav"
)

func main() {
	fs := &webdav.Handler{
		FileSystem: webdav.Dir("."),
		LockSystem: webdav.NewMemLS(),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" && handleDirList(fs.FileSystem, w, req) {
			return
		}
		// switch req.Method {
		// case "PUt", "DELETE", "PROPATCH", "MKCOL", "COPY", "MOVE":
		// 	http.Error(w, "WebDav: Read Only!!!", http.StatusForbidden)
		// 	return
		// }
		fs.ServeHTTP(w, req)
	})

	http.ListenAndServe(":8080", nil)
}

func handleDirList(fs webdav.FileSystem, w http.ResponseWriter, req *http.Request) bool {
	ctx := context.Background()
	f, err := fs.OpenFile(ctx, req.URL.Path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}

	defer f.Close()

	if fi, _ := f.Stat(); fi != nil && !fi.IsDir() {
		return false
	}

	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print(w, "Error Reading Directory!", http.StatusInternalServerError)
		return false
	}

	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	fmt.Fprintf(w, "<pre>\n")
	for _, d := range dirs {
		name := d.Name()
		if d.IsDir() {
			name += "/"
		}
		fmt.Fprintf(w, "<a href=\"%s\">%s</a>\n", name, name)
	}

	fmt.Fprintf(w, "</pre>\n")

	return true
}
