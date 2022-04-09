// Gwiki: main
// ©José Bonnet
package main

import (
  "context"
  "database/sql"
  "flag"
  "fmt"
  "log"
  "net/http"
  "os"
  "time"

  _ "github.com/lib/pq"
)

// Concentrate all app dependencies
type application struct {
  errorLog *log.Logger
  infoLog  *log.Logger
}

type config struct {
  port int
  env  string
  db   struct {
    user         string
    password     string
    host         string
    name         string
    params       string
    maxOpenConns int
    maxIdleConns int
    maxIdleTime  string
  }
}

// main glues all pieces together and starts the server
func main(){

  // Process config flags
  var cfg config

  flag.IntVar(&cfg.port, "port", 4000, "HTTP network port")
  flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
  flag.StringVar(&cfg.db.user, "db-user", os.Getenv("GWIKI_DB_USER"), "PostgreSQL database user")
  flag.StringVar(&cfg.db.password, "db-password", os.Getenv("GWIKI_DB_PASSWORD"), "PostgreSQL database password")
  flag.StringVar(&cfg.db.name, "db-name", os.Getenv("GWIKI_DB_NAME"), "PostgreSQL database name")
  flag.StringVar(&cfg.db.host, "db-host", os.Getenv("GWIKI_DB_HOST"), "PostgreSQL database host")
  flag.StringVar(&cfg.db.user, "db-params", "", "PostgreSQL database params (param1=X&param2[]=Y)")
  flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
  flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
  flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
  flag.Parse()

  // Separate INFO and ERROR logging
  infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
  errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

  // New object containing dependencies
  app := &application{
    errorLog: errorLog,
    infoLog: infoLog,
  }

  db, err := openDB(cfg)
  if err != nil {
    errorLog.Fatal(err)
  }
  defer db.Close()
  infoLog.Println("database connection pool established...")

  // Merge all configs into one single object, allowing for HTTP logging to go to the same log as the rest
  srv := &http.Server{
    Addr:     fmt.Sprintf(":%d", cfg.port),
    ErrorLog: errorLog,
    Handler:  app.routes(),
  }

  infoLog.Printf("Starting server on %s...", srv.Addr)
  err = srv.ListenAndServe()
  errorLog.Fatal(err)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
  // Use sql.Open() to create an empty connection pool, using the DSN from the config
  // struct.
  var dbDSN string

  if cfg.db.params == "" {
    dbDSN = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.db.user, cfg.db.password, cfg.db.host, cfg.db.name)
  } else {
    dbDSN = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&%s", cfg.db.user, cfg.db.password, cfg.db.host, cfg.db.name, cfg.db.params)
  }

  db, err := sql.Open("postgres", dbDSN)
  if err != nil {
    return nil, err
  }

  db.SetMaxOpenConns(cfg.db.maxOpenConns)
  db.SetMaxIdleConns(cfg.db.maxIdleConns)
  duration, err := time.ParseDuration(cfg.db.maxIdleTime)
  if err != nil {
    return nil, err
  }

  // Set the maximum idle timeout.
  db.SetConnMaxIdleTime(duration)

  // Create a context with a 5-second timeout deadline.
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  // Use PingContext() to establish a new connection to the database, passing in the
  // context we created above as a parameter. If the connection couldn't be
  // established successfully within the 5 second deadline, then this will return an
  // error.
  err = db.PingContext(ctx)
  if err != nil {
    return nil, err
  }

  // Return the sql.DB connection pool.
  return db, nil
}
