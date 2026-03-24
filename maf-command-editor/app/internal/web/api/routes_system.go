package api

import "net/http"

func (a apiRouter) registerSystemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/items", http.StatusFound)
	})
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	})
	mux.HandleFunc("POST /api/save", func(w http.ResponseWriter, r *http.Request) {
		result := a.service.ExportDatapack()
		status := http.StatusOK
		if !result.OK {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, result)
	})
}
