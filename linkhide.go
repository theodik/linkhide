package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "sync"
  "time"

  "./shortener"
)

// BaseURL is base url for link generating; defaults to localhost
// Can be set through BASE_URL env var
var BaseURL = "http://localhost"

type link struct {
  Original  string
  Short     string
  CreatedAt time.Time
}

type database struct {
  lock   sync.RWMutex
  shorts map[string]*link
}

var links database

func initializeLinks() {
  links = database{}
  links.shorts = make(map[string]*link)
}

func (db *database) Get(short string) (*link, bool) {
  db.lock.RLock()
  defer db.lock.RUnlock()
  link, ok := db.shorts[short]
  return link, ok
}

func (db *database) Set(url string) *link {
  db.lock.RLock()
  defer db.lock.RUnlock()
  short := shortener.Encode(url)
  if link, exists := db.shorts[short]; exists {
    return link
  }
  newLink := &link{Original: url, Short: short, CreatedAt: time.Now()}
  db.shorts[short] = newLink
  return newLink
}

func write404(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, "404 Document not found :)")
}

func create(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    w.WriteHeader(http.StatusNotFound)
    write404(w, r)
    return
  }

  var message struct {
    URL string `json:"url"`
  }
  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&message); err != nil || len(message.URL) > 2000 {
    log.Printf("POST /create [FAIL] %s %s", message.URL, err)
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Fprint(w, "{ \"error\": \"Could not create link.\"}")
    return
  }

  newLink := links.Set(message.URL)

  w.WriteHeader(http.StatusCreated)
  log.Printf("POST /create %s ->%s", message.URL, newLink.Short)

  fmt.Fprintf(w, "{ \"shortUrl\": \"%s/%s\"}", BaseURL, newLink.Short)
}

func index(w http.ResponseWriter, r *http.Request) {
  if r.URL.RequestURI() == "/" {
    http.ServeFile(w, r, "html/index.html")
    return
  }

  if link, ok := links.Get(r.URL.RequestURI()[1:]); ok {
    log.Printf("GET %s Redirecting to %s", r.URL.RequestURI(), link.Original)
    http.Redirect(w, r, link.Original, http.StatusMovedPermanently)
    return
  }

  w.WriteHeader(http.StatusNotFound)
  fmt.Fprint(w, "404 Document not found :)")
}

func main() {
  log.Print("Started")

  if url := os.Getenv("BASE_URL"); len(url) > 0 {
    BaseURL = url
  }

  initializeLinks()

  http.HandleFunc("/create", create)
  http.HandleFunc("/", index)
  log.Fatal(http.ListenAndServe(":80", nil))
}
