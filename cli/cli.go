package cli

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/d-kuro/egmock/serve"

	"github.com/d-kuro/egmock/logger"
)

const (
	exitCodeOK = 0
	// Exclude special meaningful exit code
	exitCodeParseFlagError = iota + 64
	exitCodeInvalidArguments
	exitCodeServeError
)

func Run(args []string) int {
	flags := flag.NewFlagSet("egmock", flag.ContinueOnError)
	status := flags.Int("s", 200, "HTTP status code")
	port := flags.String("p", "8080", "listen port number")
	body := flags.String("r", "", "response body")

	if err := flags.Parse(args[1:]); err != nil {
		return exitCodeParseFlagError
	}

	if len(flags.Args()) < 1 {
		logger.ELog.Println("invalid arguments")
		return exitCodeInvalidArguments
	}
	path := flags.Arg(0)

	logger.ILog.Println("setup mock server...")

	http.Handle(path, serve.NewMock(*status, *body))

	srv := &http.Server{Addr: ":" + *port}

	exitCh := make(chan struct{})
	go func() {
		logger.ILog.Println("start mock server")
		logger.ILog.Printf("curl http://localhost:%s%s\n", *port, path)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ELog.Println("listen and serve error:", err)
			close(exitCh)
		}
	}()

	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)

	shutdown := func() {
		// new line
		fmt.Print("\n")
		logger.ILog.Println("shutdown server")
		ctx := context.Background()
		if err := srv.Shutdown(ctx); err != nil {
			logger.ELog.Println("shutdown server error:", err)
		}
	}

	select {
	case <-sigCh:
		shutdown()
		return exitCodeOK
	case <-exitCh:
		return exitCodeServeError
	}
}
