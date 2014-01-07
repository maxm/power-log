package main

type Dispatch struct {
    f       func()
    result  chan int
}

var dispatchChan = make(chan Dispatch, 100)

func dispatchAsync(f func()) {
  dispatchChan <- Dispatch{f, nil}
}

func dispatchSync(f func()) {
  c := make(chan int)
  dispatchChan <- Dispatch{f, c}
  <-c
}

func dispatchLoop() {
  for {
    d := <-dispatchChan
    d.f()
    if d.result != nil {
      d.result <- 0
    }
  }
}