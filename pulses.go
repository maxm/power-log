package main

import (
    "fmt"
    "os"
    "time"
    "github.com/maxm/power-log/proto"
    "strings"
    "io/ioutil"
    "code.google.com/p/goprotobuf/proto"
)

var logs []*powerLog.PowerLog = make([]*powerLog.PowerLog, 0, 36500)
var currentLog *powerLog.PowerLog = nil

func fileLogName(when int64) string {
    return time.Unix(0, when).UTC().Format("2006-01-02")
}

func saveCurrentLog(logPath string) {
  if currentLog == nil || len(logPath) == 0 { return }
  data, err := proto.Marshal(currentLog)
  check(err)
  fd, err := os.Create(logPath + fileLogName((*currentLog.Start) * int64(time.Millisecond)))
  check(err)
  fd.Write(data)
  fd.Close()
}

func loadLogs(logPath string) {
  if !strings.HasSuffix(logPath, "/") { logPath += "/" }
  today := fileLogName(time.Now().UTC().UnixNano());
  fs, err := ioutil.ReadDir(logPath)
  check(err)
  for _, f := range fs {
    filename := logPath + f.Name()
    fd, err := os.Open(filename);
    if err == nil {
      buffer := make([]byte, 100000)
      fd.Read(buffer)
      log := new(powerLog.PowerLog)
      proto.Unmarshal(buffer, log)
      logs = append(logs, log)
      fmt.Printf("Open existing log %v with %v pulses\n", f.Name(), len(log.Offset))
      if f.Name() == today {
        fmt.Printf("%v is the current log\n", f.Name())
        currentLog = log
      }
    } else {
      fmt.Printf("%v\n", err)
    }
  }
}

func pulse(when, delta int64, logPath string) {
  if delta > 0 {
    fmt.Printf("%v %vW\n", time.Now().Format("Jan 2, 2006 at 15:04"), int64(time.Hour) / delta)
  }
  
  when /= int64(time.Millisecond)

  today := when - when % int64(time.Hour * 24 / time.Millisecond)
  if currentLog == nil || *currentLog.Start != today {
    currentLog = new(powerLog.PowerLog)
    currentLog.Start = &today
    currentLog.Offset = make([]int32, 0, 100000)

    if len(logs) == 0 || currentLog != logs[len(logs)-1] {
      logs = append(logs, currentLog)
    }
  }

  currentLog.Offset = append(currentLog.Offset, int32(when - *currentLog.Start))
  saveCurrentLog(logPath)
}

func check(err error) {
  if err != nil {
    fmt.Printf("Error %v\n", err)
    panic(err)
  }
}

func pollPulses(logPath string) {
  fd, err :=  os.Open("/sys/class/gpio/gpio7/value")
  if err != nil {
    fmt.Printf("Error %v\n", err)
    return
  }

  var lastPulse int64 = 0
  var previous byte = '0'
  count := 0
  current := make([]byte, 1)
  ones := false
  pulseLength := 8
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

    now := time.Now().UTC().UnixNano()

    if count == pulseLength && previous == '1' {
      ones = true
    } else if count == pulseLength && previous == '0' && ones {
      var delta int64 = 0
      if lastPulse > 0 { delta = now - lastPulse }
      dispatchAsync(func() {
        pulse(now, delta, logPath)
      })
      lastPulse = now
      ones = false
    }

    delta := int64(10 * time.Millisecond)
    time.Sleep(time.Duration(delta - now % delta))
  }
}

func listPulses(from, to int64) []int64 {
  pulses := make([]int64, 0, 10000)
  if (len(logs) == 0) { return pulses }
  log := len(logs) - 1
  start := *logs[log].Start
  for start > from && log > 0 {
    log -= 1
    start = *logs[log].Start
  }
  for ;log < len(logs); log += 1 {
    start = *logs[log].Start
    for _,offset := range logs[log].Offset {
      pulse := start + int64(offset)
      if pulse >= to {
        return pulses
      }
      if pulse >= from {
        pulses = append(pulses, pulse)
      }
    }
  }
  return pulses
}
