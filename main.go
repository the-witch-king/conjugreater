package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

//go:embed web/build/*
var buildFS embed.FS

func main() {
	static, err := fs.Sub(buildFS, "web/build")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	addr := listener.Addr().String()
	url := "http://" + addr
	fmt.Printf("Serving at %s\n", url)

	openBrowser(url)

	http.Handle("/", http.FileServer(http.FS(static)))
	log.Fatal(http.Serve(listener, nil))
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	}
	if cmd != nil {
		cmd.Start()
	}
}
