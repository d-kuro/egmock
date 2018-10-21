package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	status  int
	resBody string
}

type Request struct {
	Protocol    string `json:"protocol"`
	ContentType string `json:"content_type"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Query       string `json:"query"`
	Body        string `json:"body"`
}

func (m *mock) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// request logging
	go func() {
		bufBody := new(bytes.Buffer)
		bufBody.ReadFrom(r.Body)

		request := Request{
			Protocol:    r.Proto,
			ContentType: r.Header.Get("Content-Type"),
			Method:      r.Method,
			Path:        r.URL.Path,
			Query:       r.URL.RawQuery,
			Body:        bufBody.String(),
		}
		jsonBytes, err := json.Marshal(request)
		if err != nil {
			eLog.Println("json marshal error:", err)
		}
		iLog.Println(string(jsonBytes))
	}()

	w.WriteHeader(m.status)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(m.resBody))
	return
}

func init() {
	iLog = log.New(os.Stdout, "[info]", log.LstdFlags|log.LUTC)
	eLog = log.New(os.Stderr, "[error]", log.LstdFlags|log.LUTC)
}

func main() {
	status := flag.Int("s", 200, "HTTP status code")
	port := flag.String("p", "8080", "listen port number")
	body := flag.String("d", "", "response body")
	flag.Parse()

	if len(flag.Args()) < 1 {
		eLog.Println("invalid arguments")
	}
	path := flag.Arg(0)

	http.Handle(path, &mock{
		status:  *status,
		resBody: *body,
	})

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
		iLog.Println("start mock server")
		iLog.Printf("curl http://localhost:%s%s\n", *port, path)
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
