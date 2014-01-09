package main

import (
    "net/http"
    "strconv"
    "math"
    "errors"
    "encoding/json"
    "fmt"
)

type RangeResult struct {
  Pulses []int64
}

func parseInt(s []string) (int64, error) {
  if len(s) == 0 {
    return 0, errors.New("")
  }
  i,err := strconv.ParseInt(s[0], 10, 64)
  return i, err
}

func httpServer() {
  http.HandleFunc("/range", func(w http.ResponseWriter, r *http.Request) {
    result := RangeResult{make([]int64,0)}
    r.ParseForm()
    from, err := parseInt(r.Form["from"])
    if err == nil {
      to, err := parseInt(r.Form["to"])
      if err != nil { to = math.MaxInt64}
      dispatchSync(func() {
        result.Pulses = listPulses(from, to)
      })
    } else {
      fmt.Printf("Range request with invalid parameters %v\n", r.Form)
    }
    data, _ := json.Marshal(result)
    w.Header()["Access-Control-Allow-Origin"] = []string{"*"}
    w.Write(data)
  })
  http.Handle("/", http.FileServer(http.Dir("web")))
  http.ListenAndServe(":8080", nil)
}
