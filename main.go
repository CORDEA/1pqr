package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/qpliu/qrencode-go/qrencode"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
	"text/template"
	"time"
)

type server struct {
	tmpl   *template.Template
	param  *renderParam
	server *http.Server
}

type renderParam struct {
	Value string
	Image string
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Detected access!")
	fmt.Println(r.Header.Get("User-Agent"))
	if err := s.tmpl.Execute(w, s.param); err != nil {
		log.Fatalln(err)
	}
}

func (s *server) serve() {
	http.HandleFunc("/", s.handler)
	s.server = &http.Server{Addr: ":8080", Handler: nil}
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}

func (s *server) shutdown() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Requires name or ID.")
	}
	arg := os.Args[1]
	cmd := exec.Command("op", "item", "get", arg, "--fields", "label=password")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}
	pass := out.String()
	qr, err := qrencode.Encode(pass, qrencode.ECLevelQ)
	if err != nil {
		log.Fatalln(err)
	}

	rawImg := qr.Image(8)
	var img bytes.Buffer
	if err := png.Encode(&img, rawImg); err != nil {
		log.Fatalln(err)
	}
	b64 := base64.StdEncoding.EncodeToString(img.Bytes())

	t := "template.html"
	tmpl, err := template.New(t).ParseFiles(t)
	if err != nil {
		log.Fatalln(err)
	}

	s := server{
		tmpl: tmpl,
		param: &renderParam{
			Value: pass,
			Image: b64,
		},
	}
	fmt.Println("Starting the web server...")
	timer := time.NewTimer(3 * time.Minute)
	go func() {
		<-timer.C
		fmt.Println("Closing the web server...")
		s.shutdown()
	}()
	s.serve()
}
