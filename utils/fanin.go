package utils

import (
    "context"
    "sync"
)

// FanIn multiplexes supplied channels into one channel
func FanIn(
    ctx context.Context,
    channels ...<-chan interface{},
) <-chan interface{} {
    var wg sync.WaitGroup
    multiplexedStream := make(chan interface{})

    multiplex := func(c <-chan interface{}) {
        defer wg.Done()
        for i := range c {
            select {
            case <-ctx.Done():
                return
            case multiplexedStream <- i:
            }
        }
    }

    // Select from all the channels
    wg.Add(len(channels))
    for _, c := range channels {
        go multiplex(c)
    }

    // Wait for all the reads to complete
    go func() {
        wg.Wait()
        close(multiplexedStream)
    }()

    return multiplexedStream
}

// FanInChannelFunc turns funcs returning channel into one single channel
func FanInChannelFunc(
    ctx context.Context,
    channelFuncs ...func() <-chan interface{},
) <-chan interface{} {
    var wg sync.WaitGroup
    multiplexedStream := make(chan interface{})

    multiplex := func(c <-chan interface{}) {
        defer wg.Done()
        for i := range c {
            select {
            case <-ctx.Done():
                return
            case multiplexedStream <- i:
            }
        }
    }

    // Select from all the channelFuncs
    wg.Add(len(channelFuncs))
    for _, c := range channelFuncs {
        go multiplex(c())
    }

    // Wait for all the reads to complete
    go func() {
        wg.Wait()
        close(multiplexedStream)
    }()

    return multiplexedStream
}

// FanInFunc turns ordinary funcs into one single channel
func FanInFunc(
    ctx context.Context,
    channelFuncs ...func() interface{},
) <-chan interface{} {
    var wg sync.WaitGroup
    multiplexedStream := make(chan interface{})

    multiplex := func(c <-chan interface{}) {
        defer wg.Done()
        for i := range c {
            select {
            case <-ctx.Done():
                return
            case multiplexedStream <- i:
            }
        }
    }

    // Setup standard handling with supplied func
    channelizer := func(chanFunc func() interface{}) <-chan interface{} {
        stream := make(chan interface{})
        go func() {
            defer close(stream)
            res := chanFunc()
            stream <- res
        }()
        return stream
    }

    // Select from all the channelFuncs
    wg.Add(len(channelFuncs))
    for _, cfun := range channelFuncs {
        go multiplex(channelizer(cfun))
    }

    // Wait for all the reads to complete
    go func() {
        wg.Wait()
        close(multiplexedStream)
    }()

    return multiplexedStream
}
