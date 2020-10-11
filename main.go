package main

import (
	"context"
	"go-microservice-webinar/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	l := log.New(os.Stdout, "product-api: ", log.LstdFlags)

	ph := handlers.NewProducts(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	go func() {
		l.Println("Server listening on http://localhost:9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sc := make(chan os.Signal)
	signal.Notify(sc, os.Kill)
	signal.Notify(sc, os.Interrupt)

	sg := <-sc
	l.Println("Received terminate, graceful shutdown", sg)
	cn, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(cn)
}
