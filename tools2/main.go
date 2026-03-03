package main

import (
	"log"
	"net/http"

	"github.com/a-h/templ"

	"tools2/internal/form"
	"tools2/internal/validation"
	"tools2/views"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("POST /contact/submit", submitHandler)

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, views.Page(form.State{}))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	state := form.State{
		Name:  r.FormValue("name"),
		Phone: r.FormValue("phone"),
	}
	state.Errors = validation.Validate(state.Name, state.Phone)
	if !state.Errors.Any() {
		state.Success = "送信しました。"
	}

	if isHTMX(r) {
		render(w, r, views.ContactForm(state))
		return
	}
	render(w, r, views.Page(state))
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "render failed", http.StatusInternalServerError)
	}
}
