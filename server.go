package main

import (
  "fmt"
  "github.com/gorilla/mux"
  "io/ioutil"
  "log"
  "net/http"
  "os"
)

func main() {
  r := mux.NewRouter()

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }

  r.PathPrefix("/www/").Handler(http.StripPrefix("/www/", http.FileServer(http.Dir(getStaticDir()))))
  r.Path("/").HandlerFunc(handleIndex).Name("home")

  http.Handle("/", r)
  fmt.Printf("server started on port %s\n", port)
  log.Fatal(http.ListenAndServe(":"+port, nil))
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
