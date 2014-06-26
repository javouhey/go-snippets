package main

import (
    "fmt"
    "bytes"
    "io"
    "net/http"
    "errors"
    "regexp"
    "strconv"
    "strings"
    "encoding/json"

    "github.com/stephens2424/muxchain"
    "github.com/stephens2424/muxchain/muxchainutil"
)

var (
    re1         = regexp.MustCompile(`^\/resource\/\d+(\/)*$`)
    InvalidPath = errors.New("Path not compliant with /resource/{id:[0-9]+}")
)

func parseId(path string) (uint64, error) {
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

// @TODO does not work for HEAD
func main() {
    replyHandler := http.HandlerFunc(reply)
    idcheckHandler := http.HandlerFunc(idcheck)
    verbHandler := http.HandlerFunc(verbcheck)

    pathHandler := muxchainutil.NewPathMux()
    pathHandler.Handle("/resource/:id", idcheckHandler)

    muxchain.Chain("/resource/", verbHandler, pathHandler, replyHandler)
    fmt.Println("listening on port 9999 --> /resource/{id:[0-9]+}")
    http.ListenAndServe(":9999", muxchain.Default)
}


func idcheck(w http.ResponseWriter, req *http.Request) {
    if _, err := parseId(req.URL.Path); err != nil {
        http.Error(w, ":id must be numeric", http.StatusBadRequest)
    }
}

// w is actually a muxchain.checked which wraps a http.ResponseWriter
func verbcheck(w http.ResponseWriter, req *http.Request) {
    fmt.Println("method = ", req.Method)
    if req.Method != "GET" {
        http.Error(w, "Only GET verb is accepted", http.StatusMethodNotAllowed)
        // w.(http.Flusher).Flush() <- this is not a Flusheable ResponseWriter
    }
    fmt.Println("returned from verbcheck")
}

func reply(w http.ResponseWriter, req *http.Request) {
    id, _ := parseId(req.URL.Path)
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
