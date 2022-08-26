package WebHook

import (
    "sync"
    "time"
)

type (
    T      string
    buffer struct {
        mutex   sync.RWMutex
        buffer  []T
        buffers [][]T
        opt     *option
    }
    option struct {
        maxBuffers    uint
        fireThreshold int
        fireInterval  time.Duration
        fn            func([]T)
    }
    Option func(*option)
)

func NewBuffer(opts ...Option) *buffer {
    opt := &option{
        maxBuffers:    50,
        fireThreshold: 10,
        fireInterval:  time.Second * 5,
        fn:            func([]T) {},
    }
    for _, fn := range opts {
        fn(opt)
    }
    buf := &buffer{
        opt: opt,
    }
    buf.buffer = make([]T, 0, opt.fireThreshold)
    buf.buffers = make([][]T, opt.maxBuffers)
    go func() {
        for {
            time.Sleep(opt.fireInterval)
            if buf.Count() == 0 {
                continue
            }
            buf.flush()
        }
    }()
    return buf
}
func (b *buffer) Add(s T) {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    b.buffer = append(b.buffer, s)
    if len(b.buffer) >= b.opt.fireThreshold {
        oldBuffer := b.buffer
        b.buffer = b.buffers[0]
        b.buffers = b.buffers[1:]
        if len(b.buffers) == 0 {
            b.buffers = make([][]T, b.opt.maxBuffers)
        }
        b.opt.fn(oldBuffer)
    }
}
func (b *buffer) flush() {
    b.mutex.Lock()
    defer b.mutex.Unlock()

    oldBuffer := b.buffer
    b.buffer = b.buffers[0]
    b.buffers = b.buffers[1:]
    if len(b.buffers) == 0 {
        b.buffers = make([][]T, b.opt.maxBuffers)
    }
    b.opt.fn(oldBuffer)
}
func (b *buffer) Count() int {
    b.mutex.RLock()
    defer b.mutex.RUnlock()
    return len(b.buffer)
}
func (b *buffer) Drain() {
    if b.Count() == 0 {
        return
    }
    b.mutex.Lock()
    defer b.mutex.Unlock()
    b.opt.fn(b.buffer)
}

func WithHandler(fn func([]T)) Option {
    return func(o *option) {
        o.fn = fn
    }
}
func WithBufferSize(n uint) Option {
    return func(o *option) {
        o.maxBuffers = n
    }
}
func WithFireThreshold(n int) Option {
    return func(o *option) {
        o.fireThreshold = n
    }
}
func WithFireInterval(duration time.Duration) Option {
    return func(o *option) {
        o.fireInterval = duration
    }
}
