package ActionLog

import (
    "fmt"
    "testing"
    "time"
)

func BenchmarkBuffer(b *testing.B) {
    b.ReportAllocs()

    buf := NewRotateBuffer()
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
        sumT,
        sumItem,
    )
}
