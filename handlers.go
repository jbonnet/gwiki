package main

import (
  "fmt"
  "net/http"
)

// home is a handler that simply outputs a relevant sentence
func home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }
  w.Write([]byte("Hello from GWiki"))
}

// Add a showPage handler function.
func showPage(w http.ResponseWriter, r *http.Request) {
  slug := r.URL.Query().Get("slug")
  if slug == "" {
    http.NotFound(w, r)
    return
  }
  fmt.Fprintf(w, "Display a specific page with slug '%s'...", slug)
}

// Add a createPage handler function.
func createPage(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", http.MethodPost)
    // w.WriteHeader(405)
    // w.Write([]byte("Method Not Allowed"))
    http.Error(w, "Method Not Allowed", 405)
    return
  }
  w.Write([]byte("Create a new page..."))
}

