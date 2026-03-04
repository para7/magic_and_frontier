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
	mux.HandleFunc("GET /contact/mode-fields", modeFieldsHandler)
	mux.HandleFunc("POST /contact/submit", submitHandler)

	addr := ":8080"
	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render(w, r, views.Page(form.State{Mode: form.ModeLatLng}))
}

func modeFieldsHandler(w http.ResponseWriter, r *http.Request) {
	state := readStateFromRequest(r)
	render(w, r, views.ModeFields(state))
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	state := readStateFromRequest(r)
	state.Errors = validation.Validate(state)
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

func readStateFromRequest(r *http.Request) form.State {
	state := form.State{
		Name:      r.FormValue("name"),
		Phone:     r.FormValue("phone"),
		Mode:      form.NormalizeMode(r.FormValue("mode")),
		Latitude:  r.FormValue("latitude"),
		Longitude: r.FormValue("longitude"),
		Birthdate: r.FormValue("birthdate"),
		Height:    r.FormValue("height"),
		Weight:    r.FormValue("weight"),
	}

	switch state.Mode {
	case form.ModeLatLng:
		state.Birthdate = ""
		state.Height = ""
		state.Weight = ""
	case form.ModeBirthdate:
		state.Latitude = ""
		state.Longitude = ""
		state.Height = ""
		state.Weight = ""
	case form.ModeHeightWeight:
		state.Latitude = ""
		state.Longitude = ""
		state.Birthdate = ""
	}

	return state
}

func render(w http.ResponseWriter, r *http.Request, component templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		http.Error(w, "render failed", http.StatusInternalServerError)
	}
}
