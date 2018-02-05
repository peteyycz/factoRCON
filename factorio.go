package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"

	rcon "github.com/gtaylor/factorio-rcon"
	"golang.org/x/crypto/acme/autocert"
)

var (
	rconServerAddress  = os.Getenv("RCON_SERVER_ADDRESS")
	rconServerPassword = os.Getenv("RCON_SERVER_PASSWORD")
	serverAddress      = os.Getenv("SERVER_ADDRESS")
)

func main() {
	log.SetFlags(0)

	rconClient, err := rcon.Dial(rconServerAddress)
	if err != nil {
		panic(err)
	}
	err = rconClient.Authenticate(rconServerPassword)
	if err != nil {
		panic(err)
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			err = rconClient.Close()
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/command", &RCONHandler{client: rconClient})

	if os.Getenv("environment") == "production" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache("certs"),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("peterczibik.name"),
		}

		s := &http.Server{
			Addr: "peterczibik.name",
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}

		s.ListenAndServeTLS("", "")
	} else {
		s := &http.Server{
			Addr: "127.0.0.1:8080",
		}
		s.ListenAndServe()
	}
}
