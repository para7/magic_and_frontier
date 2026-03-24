package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"tools2/app/internal/cli"
	"tools2/app/internal/config"
	"tools2/app/internal/httpapi"
)

func main() {
	cfg := config.Load()
	args := os.Args[1:]
	if len(args) == 0 {
		os.Exit(runEditor(nil, cfg))
	}

	switch args[0] {
	case "editor":
		os.Exit(runEditor(args[1:], cfg))
	default:
		os.Exit(cli.Run(args, os.Stdout, os.Stderr, cfg))
	}
}

func runEditor(args []string, cfg config.Config) int {
	fs := flag.NewFlagSet("editor", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected args for editor: %s\n", strings.Join(fs.Args(), " "))
		return 2
	}

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, newHandler()); err != nil {
		fmt.Fprintf(os.Stderr, "editor failed: %v\n", err)
		return 1
	}
	return 0
}

func newHandler() http.Handler {
	cfg := config.Load()
	return chain(httpapi.NewHandler(cfg, httpapi.DefaultDependencies(cfg)), recoverMiddleware, loggingMiddleware)
}

type middleware func(http.Handler) http.Handler

func chain(next http.Handler, middlewares ...middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		next = middlewares[i](next)
	}
	return next
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			log.Printf("panic recovered: %v\n%s", rec, debug.Stack())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}()

		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, rec.statusCode, time.Since(start))
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
