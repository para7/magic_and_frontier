package main

import (
	"net/http"

	"tools2/app/internal/forms"
	"tools2/app/views/form"
)

func handleForm(w http.ResponseWriter, r *http.Request) {
	fs := forms.FromQuery(r)
	errs := forms.FieldErrors{}
	if err := form.Page(fs, errs, false).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleFormValidate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fs := forms.FromPostForm(r)
	errs := forms.Validate(fs)
	if err := form.FormContent(fs, errs).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleFormSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fs := forms.FromPostForm(r)
	errs := forms.Validate(fs)
	if errs.Any() {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}
	if !errs.Any() {
		w.Header().Set("HX-Trigger", "保存しました")
	}
	if err := form.Page(fs, errs, !errs.Any()).Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
