package main

import (
  "flag"
  "log"
  "net/http"
  "os"
)

func main(){
  // addr is a flag allowing different HTTP addresses
  addr := flag.String("addr", ":4000", "HTTP network address")
  flag.Parse()

  // Separate INFO and ERROR logging
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  mux := http.NewServeMux()
  mux.HandleFunc("/", home)
  mux.HandleFunc("/pages", showPage)
  mux.HandleFunc("/pages/create", createPage)

  // Serve static UI files
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  infoLog.Printf("Starting server on %s...", *addr)
  err := http.ListenAndServe(*addr, mux)
  errorLog.Fatal(err)
}
