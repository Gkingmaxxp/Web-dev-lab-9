package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/tickets/create", http.StatusSeeOther)
    })
    mux.HandleFunc("GET /tickets/create", ticketCreate)
    mux.HandleFunc("POST /tickets/create", ticketCreatePost)
    mux.HandleFunc("GET /tickets", listTickets)
    mux.HandleFunc("GET /software-requests/create", softwareCreate)
    mux.HandleFunc("POST /software-requests/create", softwareCreatePost)
    mux.HandleFunc("GET /software-requests", listSoftwareRequests)
    mux.Handle("GET /static/", http.StripPrefix("/static/",
        http.FileServer(http.Dir("static"))))
    log.Print("Open http://localhost:4000/tickets/create")
    log.Fatal(http.ListenAndServe(":4000", mux))
}
