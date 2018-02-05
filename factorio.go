package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"

	rcon "github.com/gtaylor/factorio-rcon"
	"golang.org/x/crypto/acme/autocert"

	"github.com/kelseyhightower/envconfig"
)

type specification struct {
	RconServerAddress  string   `split_words:"true"`
	RconServerPassword string   `split_words:"true"`
	ServerAddress      string   `split_words:"true"`
	CertHostWhitelist  []string `split_words:"true"`
}

func main() {
	log.SetFlags(0)
	var s specification
	err := envconfig.Process("factorcon", &s)

	rconClient, err := rcon.Dial(s.RconServerAddress)
	if err != nil {
		panic(err)
	}
	err = rconClient.Authenticate(s.RconServerPassword)
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
			HostPolicy: autocert.HostWhitelist(s.CertHostWhitelist...),
		}

		s := &http.Server{
			Addr: s.ServerAddress,
			TLSConfig: &tls.Config{
				GetCertificate: m.GetCertificate,
			},
		}

		s.ListenAndServeTLS("", "")
	} else {
		s := &http.Server{
			Addr: s.ServerAddress,
		}
		s.ListenAndServe()
	}
}
