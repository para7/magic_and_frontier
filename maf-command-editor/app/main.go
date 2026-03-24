package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"tools2/app/internal/application"
	"tools2/app/internal/config"
	"tools2/app/internal/web/api"
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
	case "validate":
		os.Exit(runValidate(args[1:], os.Stdout, os.Stderr, cfg))
	case "export":
		os.Exit(runExport(args[1:], os.Stdout, os.Stderr, cfg))
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", args[0])
		printUsage(os.Stderr)
		os.Exit(2)
	}
}

func runValidate(args []string, stdout, stderr io.Writer, cfg config.Config) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc := application.NewService(cfg, application.DefaultDependencies(cfg))
	report, err := svc.ValidateAll()
	if err != nil {
		fmt.Fprintf(stderr, "validate failed: %v\n", err)
		return 1
	}
	if !report.OK {
		fmt.Fprintln(stderr, "savedata validation failed:")
		fmt.Fprintln(stderr, report.String())
		return 1
	}

	fmt.Fprintf(stdout, "savedata validation ok: items=%d grimoire=%d skills=%d enemy_skills=%d enemies=%d spawn_tables=%d treasures=%d loottables=%d\n",
		report.Counts.Items,
		report.Counts.Grimoire,
		report.Counts.Skills,
		report.Counts.EnemySkills,
		report.Counts.Enemies,
		report.Counts.SpawnTables,
		report.Counts.Treasures,
		report.Counts.LootTables,
	)
	return 0
}

func runExport(args []string, stdout, stderr io.Writer, cfg config.Config) int {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	fs.SetOutput(stderr)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc := application.NewService(cfg, application.DefaultDependencies(cfg))
	result := svc.ExportDatapack()
	if !result.OK {
		fmt.Fprintf(stderr, "export failed: %s\n", result.Message)
		if result.Details != "" {
			fmt.Fprintln(stderr, result.Details)
		}
		return 1
	}

	fmt.Fprintf(stdout, "datapack export completed: output=%s total_files=%d item_functions=%d item_loot=%d spell_functions=%d spell_loot=%d skills=%d enemy_skills=%d enemy_functions=%d enemy_loot=%d treasure_loot=%d loottable_loot=%d\n",
		result.OutputRoot,
		result.Generated.TotalFiles,
		result.Generated.ItemFunctions,
		result.Generated.ItemLootTables,
		result.Generated.SpellFunctions,
		result.Generated.SpellLootTables,
		result.Generated.SkillFunctions,
		result.Generated.EnemySkillFunctions,
		result.Generated.EnemyFunctions,
		result.Generated.EnemyLootTables,
		result.Generated.TreasureLootTables,
		result.Generated.LoottableLootTables,
	)
	return 0
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
	return chain(api.NewHandler(cfg, application.Dependencies{}), recoverMiddleware, loggingMiddleware)
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
		if strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}
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

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "usage: mce <command>")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "commands:")
	fmt.Fprintln(w, "  editor     start web editor server")
	fmt.Fprintln(w, "  validate   validate savedata and export settings")
	fmt.Fprintln(w, "  export     validate and export datapack")
}
