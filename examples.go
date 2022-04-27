package main

import (
    "context"
    "fmt"
    "time"

    "github.com/cekrem/goutils/utils"
)

func main() {
    ctx := context.Background()
    res := utils.FanInChannelFunc(ctx, func() <-chan interface{} {
        stream := make(chan interface{})
        go func() {
            defer close(stream)
            time.Sleep(1 * time.Second)
            stream <- "foo"
        }()
        return stream
    }, func() <-chan interface{} {
        stream := make(chan interface{})
        go func() {
            defer close(stream)
            time.Sleep(1 * time.Second)
            stream <- "bar"
        }()
        return stream
    })

    for entry := range res {
        fmt.Println(entry)
    }

    ctx, cancel := context.WithCancel(ctx)
    cancel()

    res = utils.FanInFunc(ctx, func() interface{} {
        time.Sleep(100)
        return "this should never return: canceled!"
    })

    fmt.Println(<-res)
}
