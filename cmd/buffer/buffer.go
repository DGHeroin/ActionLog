package main

import (
    "fmt"
    "github.com/DGHeroin/ActionLog"
    "sync"
    "time"
)

func main() {
    t0()
   // t1()
}
func t0() {
    buf := ActionLog.NewRotateBuffer()
    buf.EnableGzip(false)
    sumSize := 0
    sumT := time.Duration(0)
    sumItem := 0
    sumRaw := 0
    buf.OnRotate(func(data []byte, t0, t1 time.Time, num int) {
        dt := t1.Sub(t0)
        raw, _ := buf.Decompress(data)
        sz := len(data)
        szRaw := len(raw)
        fmt.Println("rotate size:", ActionLog.HumanFileSize(float64(sz)), ActionLog.HumanFileSize(float64(szRaw)), dt, num)
        sumSize += sz
        sumT += dt
        sumItem += num
        sumRaw += szRaw
    })
    
    L := ActionLog.New()
    L.SetWriter(buf)
    defer buf.Flush()
    
    startTime := time.Now()
    for i := 0; i < 1000*1000; i++ {
        L.Info(ActionLog.F{"code": 10, "playerId": "aaa"}, "hello \n", "world")
    }
    buf.Flush()
    
    fmt.Println("elapsed time:", time.Since(startTime),
        ActionLog.HumanFileSize(float64(sumSize)),
        ActionLog.HumanFileSize(float64(sumRaw)),
        sumT,
        sumItem,
    )
}

func t1() {
    buf := ActionLog.NewRotateBuffer()
    //buf.EnableGzip(false)
    sumSize := 0
    sumT := time.Duration(0)
    sumItem := 0
    sumRaw := 0
    buf.OnRotate(func(data []byte, t0, t1 time.Time, num int) {
        dt := t1.Sub(t0)
        raw, _ := buf.Decompress(data)
        sz := len(data)
        szRaw := len(raw)
        fmt.Println("rotate size:", ActionLog.HumanFileSize(float64(sz)), ActionLog.HumanFileSize(float64(szRaw)), dt, num)
        sumSize += sz
        sumT += dt
        sumItem += num
        sumRaw += szRaw
    })
    
    L := ActionLog.New()
    L.SetWriter(buf)
    defer buf.Flush()
    
    startTime := time.Now()
    wg := sync.WaitGroup{}
    fn := func() {
        defer wg.Done()
        for i := 0; i < 1000; i++ {
            L.Info(ActionLog.F{"code": 10, "playerId": "aaa"}, "hello \n", "world")
        }
    }
    wg.Add(1000)
    for i := 0; i < 1000; i++ {
        fn()
    }
    wg.Wait()
    buf.Flush()
    
    fmt.Println("elapsed time:", time.Since(startTime),
        ActionLog.HumanFileSize(float64(sumSize)),
        ActionLog.HumanFileSize(float64(sumRaw)),
        sumT,
        sumItem,
    )
}
