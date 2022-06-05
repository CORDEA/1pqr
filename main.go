package main

import (
	"bytes"
	"fmt"
	"github.com/qpliu/qrencode-go/qrencode"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Detected access!")
	fmt.Println(r.Header.Get("User-Agent"))
	if _, err := fmt.Fprint(w, "<html></html>"); err != nil {
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
	fmt.Println(qr.String())

	fmt.Println("Starting the web server...")
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
