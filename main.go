package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	iLog *log.Logger
	eLog *log.Logger
)

type mock struct {
	status int
}

func (m *mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(m.status)
	return
}

func init() {
	iLog = log.New(os.Stdout, "[info]", log.LstdFlags|log.LUTC)
	eLog = log.New(os.Stderr, "[error]", log.LstdFlags|log.LUTC)
}

func main() {
	port := flag.String("p", "8080", "Listen port number")
	status := flag.Int("s", 200, "HTTP status code")
	flag.Parse()

	if len(flag.Args()) < 1 {
		eLog.Println("invalid arguments")
	}
	path := flag.Arg(0)

	http.Handle(path, &mock{status: *status})

	srv := &http.Server{Addr: ":" + *port}
	defer func() {
		// new line
		fmt.Print("\n")
		iLog.Println("shutdown server")
		ctx := context.Background()
		if err := srv.Shutdown(ctx); err != nil {
			eLog.Println("shutdown server error:", err)
		}
	}()

	go func() {
		iLog.Printf("curl -i http://localhost:%s%s\n", *port, path)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			eLog.Println("listen and serve error:", err)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	select {
	case <-sigCh:
		// execute defer func
		return
	}

	return
}
