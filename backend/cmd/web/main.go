package main

import (
  "log"
  "net/http"
)

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", home)

  mux.HandleFunc("/snippet", showPageSnippet)
  mux.HandleFunc("/snippet/create", createPageSnippet)

  log.Println("Starting server on :4000")
  err := http.ListenAndServe(":4000", mux)
  log.Fatal(err)
}
