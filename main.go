package main

import (
    "flag"
)

func main() {
  flag.Parse()
  args := flag.Args()
  logPath := ""
  if len(args) > 0 {
    logPath = args[0]
    loadLogs(logPath)
  }

  go httpServer()
  go pollPulses(logPath)
  dispatchLoop()
}
