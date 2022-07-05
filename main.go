package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nicholasjackson/env"
	"gitlab.com/my-whoosh/admin/handlers"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for server")

func main() {
	env.Parse()

	l := log.New(os.Stdout, "game-admin-api-", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()

	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
		ReadTimeout:  1 * time.Second,   // max time to read request from clinet
		WriteTimeout: 1 * time.Second,   // max time to write response to the client
	}

	go func() {
		l.Println("starting server on port :9090")

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// create a new channel
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 1*time.Second)
	s.Shutdown(tc)
}
