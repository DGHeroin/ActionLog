package ActionLog

import (
    "bytes"
    "sync"
    "sync/atomic"
    "time"
)

type (
    RotateBuffer struct {
        mu       sync.Mutex
        buf      *bytes.Buffer
        MaxSize  int
        onRotate func(buffer []byte, t0, t1 time.Time, num int)
        onWrite  func([]byte)
        t0       time.Time
        c        int32
        wg       *sync.WaitGroup
    }
)

func NewRotateBuffer() *RotateBuffer {
    return &RotateBuffer{
        buf:     &bytes.Buffer{},
        MaxSize: 100 * 1000 * 1000, // 100M
        wg:      &sync.WaitGroup{},
    }
}

func (b *RotateBuffer) Write(p []byte) (n int, err error) {
    var isOverLimit bool

    b.mu.Lock()
    isOverLimit = b.buf.Len()+len(p) > b.MaxSize
    b.mu.Unlock()

    if isOverLimit {
        b.rotate()
    }

    num := atomic.AddInt32(&b.c, 1)
    if num == 1 {
        b.t0 = time.Now()
    }
    if b.onWrite != nil {
        b.onWrite(p)
    }
    b.mu.Lock()
    defer b.mu.Unlock()
    //data := append(p, ',')
    return b.buf.Write(p)
}

func (b *RotateBuffer) rotate() {
    b.mu.Lock()
    defer b.mu.Unlock()
    if b.onRotate != nil && b.buf.Len() > 0 {
        rBuf := &bytes.Buffer{}
        rBuf.WriteString("[\n")
        data := b.buf.Bytes()
        data = data[:len(data)-2]
        rBuf.Write(data)
        rBuf.WriteString("\n]")
        b.wg.Add(1)
        num := int(b.c)
        go func() {
            defer b.wg.Done()
            b.onRotate(rBuf.Bytes(), b.t0, time.Now(), num)
        }()
    }
    b.buf = &bytes.Buffer{}
    atomic.StoreInt32(&b.c, 0)
}

func (b *RotateBuffer) Flush() {
    b.rotate()
}
func (b *RotateBuffer) Close() {
    b.rotate()
    b.wg.Wait()
}
func (b *RotateBuffer) Current() []byte {
    b.mu.Lock()
    defer b.mu.Unlock()
    buf := &bytes.Buffer{}
    buf.Write(b.buf.Bytes())
    return buf.Bytes()
}

func (b *RotateBuffer) OnRotate(fn func(buffer []byte, t0, t1 time.Time, num int)) {
    b.onRotate = fn
}
func (b *RotateBuffer) OnWrite(fn func([]byte)) {
    b.onWrite = fn
}
