package main


import (
	"bytes"
	"io"
	"log"
	"mime"
	"net"
	"net/http"
	"path/filepath"

	"github.com/zserge/webview"
)

func startServer() string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		defer ln.Close()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if len(path) > 0 && path[0] == '/' {
				path = path[1:]
			}
			if path == "" {
				path = "index.html"
			}
			if bs, err := Asset(path); err != nil {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.Header().Add("Content-Type", mime.TypeByExtension(filepath.Ext(path)))
				io.Copy(w, bytes.NewBuffer(bs))
			}
		})
		log.Fatal(http.Serve(ln, nil))
	}()
	return "http://" + ln.Addr().String()
}


func main() {
	url := startServer()
	w := webview.New(webview.Settings{
		Title:  "RK Dashboard",
		URL:    url,
		Resizable: true,
	})
	defer w.Exit()
	w.Run()
}

