package main

import (
    "bytes"
    "html/template"
    "log"
    "net/http"
)

// any allows the same renderer to receive FormData or ListData.
func renderTemplate(w http.ResponseWriter, status int,
    filename string, data any) {
    tmpl, err := template.ParseFiles("templates/" + filename)
    if err != nil {
        log.Print(err)
        http.Error(w, "Page unavailable", 500)
        return
    }
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        log.Print(err)
        http.Error(w, "Page unavailable", 500)
        return
    }
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(status)
    if _, err := buf.WriteTo(w); err != nil {
        log.Print(err)
    }
}
