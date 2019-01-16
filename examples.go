package main

import (
    "fmt"
    "time"

    "github.com/cekrem/goutils/utils"
)

func main(){
    done := make(chan interface{})
    res := utils.FanInChannelFunc(done, func() <- chan interface{} {
        stream := make(chan interface{})
        go func () {
            defer close(stream)
            time.Sleep(1 * time.Second)
            stream <- "foo"
        }()
        return stream
    }, func() <- chan interface{} {
        stream := make(chan interface{})
        go func () {
            defer close(stream)
            time.Sleep(1 * time.Second)
            stream <- "boo"
        }()
        return stream
    })

    fmt.Println(<-res)

    res = utils.FanInFunc(done, func() interface{}{
        return "foo"
    })

    fmt.Println(<-res)

}