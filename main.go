package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/qpliu/qrencode-go/qrencode"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type server struct {
	base64Qr string
}

func (s *server) buildHtml(content string) string {
	font := "<link rel=\"preconnect\" href=\"https://fonts.googleapis.com\">" +
		"<link rel=\"preconnect\" href=\"https://fonts.gstatic.com\" crossorigin>" +
		"<link href=\"https://fonts.googleapis.com/css2?family=Roboto+Mono&display=swap\" rel=\"stylesheet\">"
	style := "style=\"text-align: center;margin: 16px;\""
	return fmt.Sprintf("<html><head>%s</head><body %s>%s</body></html>", font, style, content)
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Detected access!")
	fmt.Println(r.Header.Get("User-Agent"))
	img := fmt.Sprintf("<img src=\"data:image/png;base64,%s\" />", s.base64Qr)
	html := s.buildHtml(img)
	if _, err := fmt.Fprint(w, html); err != nil {
		log.Fatalln(err)
	}
}

func (s *server) serve() {
	http.HandleFunc("/", s.handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
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
	qr, err := qrencode.Encode(out.String(), qrencode.ECLevelQ)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Starting the web server...")
	raw := qr.Image(8)
	var img bytes.Buffer
	if err := png.Encode(&img, raw); err != nil {
		log.Fatalln(err)
	}
	b64 := base64.StdEncoding.EncodeToString(img.Bytes())

	s := server{base64Qr: b64}
	s.serve()
}
