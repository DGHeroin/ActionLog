package ActionLog

import (
    "fmt"
    "testing"
    "time"
)

func BenchmarkBuffer(b *testing.B) {
    b.ReportAllocs()

    buf := NewRotateBuffer()
    buf.EnableGzip(true)
    sumSize := 0
    sumT := time.Duration(0)
    sumItem := 0
    sumRaw := 0
    buf.OnRotate(func(data []byte, t0, t1 time.Time, num int) {
        dt := t1.Sub(t0)
        //raw, _ := buf.Decompress(data)
        sz := len(data)
        //szRaw := len(raw)
        //fmt.Println("rotate size:", HumanFileSize(float64(sz)), HumanFileSize(float64(szRaw)), dt, num)
        sumSize += sz
        sumT += dt
        sumItem += num
        //sumRaw += szRaw
    })

    L := New()
    L.SetWriter(buf)
    defer buf.Flush()

    startTime := time.Now()
    for i := 0; i < b.N; i++ {
        L.Info(F{"code": 10, "playerId": "aaa"}, "hello \n", "world")
    }
    buf.Flush()
    elapsedTime := time.Since(startTime)
    fmt.Println("elapsed time:", elapsedTime,
        HumanFileSize(float64(sumSize)),
        HumanFileSize(float64(sumRaw)),
        sumT,
        sumItem,
    )
}
