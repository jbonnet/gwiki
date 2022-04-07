// GWiki: routes
// ©José Bonnet
package main

import "net/http"

func (app *application) routes() *http.ServeMux {
  mux := http.NewServeMux()
  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/pages", app.showPage)
  mux.HandleFunc("/pages/create", app.createPage)

  // Serve static UI files
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  return mux
}
