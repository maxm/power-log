package main

import (
    "fmt"
    "os"
    "time"
    "flag"
    "strings"
    "code.google.com/p/goprotobuf/proto"
    "github.com/maxm/powerLog/proto"
)

var logPath string = ""
var currentLog *powerLog.PowerLog = nil

func check(err error) {
  if err != nil {
    fmt.Printf("Error %v\n", err)
    panic(err)
  }
}

func fileLogName(when int64) string {
  return time.Unix(0, when).Format("2006-01-02")
}

func saveCurrentLog() {
  if currentLog == nil || len(logPath) == 0 { return }
  data, err := proto.Marshal(currentLog)
  check(err)
  fd, _ := os.Create(logPath + fileLogName(*currentLog.Start))
  fd.Write(data)
  fd.Close()
}

func pulse(when, delta int64) {
  if delta > 0 {
    fmt.Printf("%v %vW\n", time.Now().Format("Jan 2, 2006 at 15:04"), int64(time.Hour) / delta)
  }

  today := when - when % int64(time.Hour * 24)
  if currentLog == nil || *currentLog.Start != today {
    currentLog = new(powerLog.PowerLog)
    currentLog.Start = &today
    currentLog.Offset = make([]int32, 0, 100000)
  }

  currentLog.Offset = append(currentLog.Offset, int32(when - *currentLog.Start))
  saveCurrentLog()
}

func pollPulses() {
  fd, err :=  os.Open("/sys/class/gpio/gpio7/value")
  check(err)

  var previous byte = '0'
  count := 0
  current := make([]byte, 1)
  ones := false
  pulseLength := 8
  var lastPulse int64 = 0
  for {
    fd.Seek(0, 0)
    _, err := fd.Read(current)
    check(err)

    if current[0] == previous {
      count += 1
    } else {
      count = 1
    }
    previous = current[0]

    now := time.Now().UnixNano()

    if count == pulseLength && previous == '1' {
      ones = true
    } else if count == pulseLength && previous == '0' && ones {
      var delta int64 = 0
      if lastPulse > 0 { delta = now - lastPulse }
      pulse(now, delta)
      lastPulse = now
      ones = false
    }
    
    delta := int64(10 * time.Millisecond)
    time.Sleep(time.Duration(delta - now % delta))
  }
}

func main() {
  args := flag.Args()
  if len(args) > 0 {
    logPath = args[0]
    if !strings.HasSuffix(logPath, "/") { logPath += "/" }
    fd, err := os.Open(logPath + fileLogName(time.Now().UnixNano()))
    if err != nil {
      buffer := make([]byte, 100000)
      fd.Read(buffer)
      proto.Unmarshal(buffer, currentLog)
    }
  }
  pollPulses()
}
