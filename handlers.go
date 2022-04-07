// GWiki handlers
// @José Bonnet
package main

import (
  "fmt"
  "html/template"
  "net/http"
)

// home is a handler that simply outputs a relevant sentence
func (app *application) home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    app.notFound(w)
    return
  }

  files := []string{
    "./ui/html/home.page.tmpl",
    "./ui/html/base.layout.tmpl",
    "./ui/html/footer.partial.tmpl",
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    app.serverError(w, err)
    return
  }

  // We then use the Execute() method on the template set to write the template
  // content as the response body. The last parameter to Execute() represents any
  // dynamic data that we want to pass in, which for now we'll leave as nil.
  err = ts.Execute(w, nil)
  if err != nil {
    app.serverError(w, err)
  }
}

// Add a showPage handler function.
func (app *application) showPage(w http.ResponseWriter, r *http.Request) {
  slug := r.URL.Query().Get("slug")
  if slug == "" {
    app.notFound(w)
    return
  }
  fmt.Fprintf(w, "Display a specific page with slug '%s'...", slug)
}

// Add a createPage handler function.
func (app *application) createPage(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    w.Header().Set("Allow", http.MethodPost)
    app.clientError(w, http.StatusMethodNotAllowed)
    return
  }
  w.Write([]byte("Create a new page..."))
}

