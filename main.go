package main

import (
  "flag"
  "log"
  "net/http"
  "os"
)

// Concentrate all app dependencies
type application struct {
  errorLog *log.Logger
  infoLog  *log.Logger
}

func main(){
  // addr is a flag allowing different HTTP addresses
  addr := flag.String("addr", ":4000", "HTTP network address")
  flag.Parse()

  // Separate INFO and ERROR logging
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  // New object containing dependencies
  app := &application{
    errorLog: errorLog,
    infoLog: infoLog,
  }

  // Merge all configs into one single object, allowing for HTTP logging to go to the same log as the rest
  srv := &http.Server{
    Addr:     *addr,
    ErrorLog: errorLog,
    Handler:  app.routes(),
  }

  infoLog.Printf("Starting server on %s...", *addr)
  err := srv.ListenAndServe()
  errorLog.Fatal(err)
}
