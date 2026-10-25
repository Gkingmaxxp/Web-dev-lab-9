package main

import (
    "log"
    "net/http"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /reports/create", reportCreate)
    mux.HandleFunc("POST /reports/create", reportCreatePost)
    mux.HandleFunc("GET /reports", listReports)
    mux.Handle("GET /static/", http.StripPrefix("/static/",
        http.FileServer(http.Dir("static"))))
    log.Print("Open http://localhost:4000/reports/create")
    log.Fatal(http.ListenAndServe(":4000", mux))
}
