package main

import (
  "fmt"
  "github.com/gorilla/mux"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "database/sql"
  _ "github.com/lib/pq"
)

func main() {
  r := mux.NewRouter()

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  r.PathPrefix("/www/").Handler(http.StripPrefix("/www/", http.FileServer(http.Dir(getStaticDir()))))
  r.Path("/").HandlerFunc(handleIndex).Name("home")
  r.Path("/api/").HandlerFunc(handleSql).Name("sql")

  http.Handle("/", r)
  fmt.Printf("server started on port %s\n", port)
  log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleSql(w http.ResponseWriter, r *http.Request) {
  username := os.Getenv("db_user")
  password := os.Getenv("db_password")
  db, err := sql.Open("postgres", "postgres://" + username + ":" + password + "@homelessdb.cmcvtt7pgaun.us-east-1.rds.amazonaws.com/homelessdb")
  if err != nil {
    fmt.Printf("Failed to connect to the database: " + err.Error())
    return
  }
  rows, err := db.Query("select * from public.user")
  if err != nil {
    fmt.Printf("Failed to query the database: " + err.Error())
    return
  }
  defer rows.Close()
  for rows.Next() {
    fmt.Printf("%s", row)
    // TODO: effort. we're switching to postgrest
  }
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
  file, err := ioutil.ReadFile(getStaticDir() + "index.html")
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Printf(err.Error())
  }

  w.Header().Set("Content-Type", "text/html")
  fmt.Fprintf(w, "%s", file)
}

func getStaticDir() string {
  staticDir := os.Getenv("www")
  if staticDir == "" {
    staticDir = "www/"
  }
  return staticDir
}
