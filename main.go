package main

import (
	"net/http"
	"time"

	"gin_websocket/router"

	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {
	server := &http.Server{
		Addr:         "1234",
		Handler:      router.MainRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return server.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		// router.Log.Fatal(err)
	}
}
