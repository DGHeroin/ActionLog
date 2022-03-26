package main

import (
    "fmt"
    "github.com/DGHeroin/ActionLog"
    "io/ioutil"
    "log"
    "os"
    "sync"
    "time"
)

func main() {
    t0()
    // t1()
}
func t0() {
    buf := ActionLog.NewRotateBuffer()
    sumSize := 0
    sumT := time.Duration(0)
    sumItem := 0
    buf.OnRotate(func(data []byte, t0, t1 time.Time, num int) {
        dt := t1.Sub(t0)
        sz := len(data)
        sumSize += sz
        sumT += dt
        sumItem += num
        timefmt := "200601-02_15-04-05-Z0700"
        key := fmt.Sprintf("%v-%v_%v.json",
            t0.Format(timefmt), t1.Format(timefmt), num)
        fmt.Println(key, len(data))

        log.Println(key)

        err := ioutil.WriteFile(key, data, os.ModePerm)
        if err != nil {
           fmt.Println(err)
           return
        }

    })

    L := ActionLog.New()
    L.SetWriter(buf)
    defer buf.Flush()

    startTime := time.Now()
    for i := 1; i <= 1000*1000; i++ {
        L.Info(ActionLog.F{"code": i, "playerId": "aaa"}, "hello \n", "world")
    }
    buf.Flush()
    L.Info(ActionLog.F{"code": -1, "playerId": "aaa"}, "hello \n", "world")
    L.Info(ActionLog.F{"code": -2, "playerId": "aaa"}, "hello \n", "world")
    buf.Close()
    elapsedTime := time.Since(startTime)
    fmt.Println("elapsed time:", elapsedTime,
        ActionLog.HumanFileSize(float64(sumSize)),
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
    buf.OnRotate(func(data []byte, t0, t1 time.Time, num int) {
        dt := t1.Sub(t0)
        sz := len(data)
        sumSize += sz
        sumT += dt
        sumItem += num
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
        sumT,
        sumItem,
    )
}
