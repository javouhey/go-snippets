package main

import (
    "fmt"
    "bytes"
    "io"
    "log"
    "net/http"
    "errors"
    "regexp"
    "strconv"
    "strings"
    "encoding/json"
)

var (
    re1         = regexp.MustCompile(`^\/resource\/\d+(\/)*$`)
    InvalidPath = errors.New("Path not compliant with /resource/{id:[0-9]+}")
)

// Resource encapsulates endpoint of REST API 
// to demonstrate hawk API.
// http://127.0.0.1:9999/resource/1?b=1&a=2
// Using built-in default http muxer
type Resource struct {
}

func (r *Resource) parseId(path string) (uint64, error) {
    if !re1.MatchString(path) {
        return 0, InvalidPath
    }

    parts := strings.Split(path, "/")
    if len(parts) < 3 {
        return 0, InvalidPath
    }

    if id, err := strconv.ParseUint(parts[2], 10, 64); err != nil {
        return 0, err
    } else {
        return id, nil
    }
}

func (r *Resource) get(w http.ResponseWriter, req *http.Request) {
    var id uint64
    if res, err := r.parseId(req.URL.Path); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        io.WriteString(w, err.Error())
        return
    } else {
        id = res
    }

    q := req.URL.Query()
    q.Set("greeting", "bonjour le monde")
    q.Set("id", fmt.Sprintf("%d", id))
    if ba, err := json.Marshal(q); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        io.WriteString(w, err.Error())
    } else {
        w.Header().Set("Content-Type", "application/json")
        io.Copy(w, bytes.NewReader(ba))
    }
}

// Handle dispatching of requests
func (r *Resource) Handle(w http.ResponseWriter, req *http.Request) {
    switch req.Method {
    case "GET":
        r.get(w, req)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func main() {
    resource := &Resource{}
    http.HandleFunc("/resource/", resource.Handle)
    fmt.Println("listening on port 9999 --> /resource/{id:[0-9]+}")
    log.Fatal(http.ListenAndServe(":9999", nil))
}
