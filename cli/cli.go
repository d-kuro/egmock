package cli

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/d-kuro/egmock/serve"

	"github.com/d-kuro/egmock/log"
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
		log.Error("invalid arguments")
		return exitCodeInvalidArguments
	}
	path := flags.Arg(0)

	log.Info("setup mock server...")

	http.Handle(path, serve.NewMock(*status, *body))

	srv := &http.Server{Addr: ":" + *port}

	exitCh := make(chan struct{})
	go func() {
		log.Info("start mock server")
		log.Info("curl http://localhost:" + *port + path)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen and serve error", zap.Error(err))
			close(exitCh)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	shutdown := func() {
		// new line
		fmt.Print("\n")
		log.Info("shutdown mock server")
		ctx := context.Background()
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("shutdown server error", zap.Error(err))
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
