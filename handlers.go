package main

import (
  "fmt"
  "html/template"
  "log"
  "net/http"
)

// home is a handler that simply outputs a relevant sentence
func home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  files := []string{
    "./ui/html/home.page.tmpl",
    "./ui/html/base.layout.tmpl",
    "./ui/html/footer.partial.tmpl",
  }
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    http.Error(w, "Internal Server Error", 500)
    return
  }
  // We then use the Execute() method on the template set to write the template
  // content as the response body. The last parameter to Execute() represents any
  // dynamic data that we want to pass in, which for now we'll leave as nil.
  err = ts.Execute(w, nil)
  if err != nil {
    log.Println(err.Error())
    http.Error(w, "Internal Server Error", 500)
  }
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

