package ActionLog

import (
    "bytes"
    "compress/gzip"
    "fmt"
    "io/ioutil"
    "sync/atomic"
    "time"
)

type (
    RotateBuffer struct {
        buf        bytes.Buffer
        MaxSize    int
        enableGzip bool
        onRotate   func(buffer []byte, t0, t1 time.Time, num int)
        onWrite    func([]byte)
        t0         time.Time
        c          int32
    }
)

func NewRotateBuffer() *RotateBuffer {
    return &RotateBuffer{
        MaxSize:    10 * 1000 * 1000, // 100M
        enableGzip: true,
    }
}
func (b *RotateBuffer) Write(p []byte) (n int, err error) {
    if b.buf.Len()+len(p) > b.MaxSize {
        b.rotate()
    }
    num := atomic.AddInt32(&b.c, 1)
    if num == 1 {
        b.t0 = time.Now()
    }
    if b.onWrite != nil {
        b.onWrite(p)
    }
    return b.buf.Write(p)
}

func (b *RotateBuffer) rotate() {
    defer b.buf.Reset()
    if b.onRotate == nil {
        return
    }
    if b.buf.Len() == 0 {
        return
    }
    data := b.buf.Bytes()
    if b.enableGzip {
        compressedData, err := b.Compress(data)
        if err != nil {
            fmt.Println(err)
            return
        }
        data = compressedData
    }

    if b.onRotate != nil && len(data) > 0 {
        b.onRotate(data, b.t0, time.Now(), int(b.c))
    }
    atomic.StoreInt32(&b.c, 0)
}

func (b *RotateBuffer) Flush() {
    b.rotate()
}
func (b *RotateBuffer) Current() []byte {
    data := b.buf.Bytes()
    if b.enableGzip {
        compressedData, err := b.Compress(data)
        if err != nil {
            return nil
        }
        data = compressedData
    }
    return data
}
func (b *RotateBuffer) EnableGzip(enable bool) {
    b.enableGzip = enable
}

func (b *RotateBuffer) OnRotate(fn func(buffer []byte, t0, t1 time.Time, num int)) {
    b.onRotate = fn
}
func (b *RotateBuffer) OnWrite(fn func([]byte)) {
    b.onWrite = fn
}
func (b *RotateBuffer) Compress(data []byte) ([]byte, error) {
    var buf bytes.Buffer
    w := gzip.NewWriter(&buf)
    if _, err := w.Write(data); err != nil {
        return nil, err
    }
    err := w.Close()
    return buf.Bytes(), err
}
func (b *RotateBuffer) Decompress(data []byte) ([]byte, error) {
    r, err := gzip.NewReader(bytes.NewReader(data))
    if err != nil {
        return nil, err
    }
    if raw, err := ioutil.ReadAll(r); err != nil {
        return nil, err
    } else {
        err = r.Close()
        return raw, err
    }
}
