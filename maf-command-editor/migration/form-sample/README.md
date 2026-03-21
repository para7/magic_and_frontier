# Form Sample Archive

This directory preserves the pre-migration `templ + htmx` demo form as reference material for the frontend migration in plan 3.

- The archived files were moved out of `app/` so they are no longer part of the Go build or test graph.
- Source files were renamed with `.txt` suffixes to keep the original contents readable without compiling them.
- The server migration in plan 2 intentionally ships API-only behavior under `/health` and `/api/*`.

Use these archived files only as implementation reference when rebuilding the real SSR UI in `migration/03_frontend_migration_runbook.md`.
