package ActionLog

import (
    "bytes"
    "io"
    "os"
    "sync"
    "time"
)

type (
    ActionLog struct {
        mutex          sync.RWMutex
        pool           sync.Pool
        standardFields F
        hooks          []Hook
        writer         io.Writer
        formatter      Formatter
    }
    F map[string]interface{}
)

func New(fields ...F) *ActionLog {
    L := &ActionLog{
        writer:         os.Stdout,
        formatter:      &defaultFormatter{},
        standardFields: map[string]interface{}{},
    }
    for _, field := range fields {
        for k, v := range field {
            L.standardFields[k] = v
        }
    }
    L.pool.New = func() interface{} {
        return &Entry{
            Data:   make(F),
            Buffer: &bytes.Buffer{},
        }
    }
    return L
}
func (a *ActionLog) SetWriter(w io.Writer) {
    a.writer = w
}
func (a *ActionLog) Info(fields F, args ...interface{}) {
    entry := a.allocEntry()
    
    entry.WithTime(time.Now()).WithFields(a.standardFields).WithFields(fields).Info(args...)
    a.fire(entry)
    a.write(entry)
    
    a.freeEntry(entry)
}
func (a *ActionLog) allocEntry() *Entry {
    entry := a.pool.Get().(*Entry)
    return entry
}

func (a *ActionLog) freeEntry(entry *Entry) {
    entry.Data = map[string]interface{}{}
    entry.Buffer.Reset()
    a.pool.Put(entry)
}
func (a *ActionLog) fire(entry *Entry) {
    a.mutex.RLock()
    defer a.mutex.RUnlock()
    if len(a.hooks) == 0 {
        return
    }
    for _, hook := range a.hooks {
        err := hook.Fire(entry)
        if err != nil {
            return
        }
    }
}
func (a *ActionLog) AddHook(hook Hook) {
    if hook == nil {
        return
    }
    a.mutex.Lock()
    defer a.mutex.Unlock()
    a.hooks = append(a.hooks, hook)
}

func (a *ActionLog) write(entry *Entry) {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    data, err := a.formatter.Format(entry)
    if err != nil {
        return
    }
    _, err = a.writer.Write(data)
    if err != nil {
        return
    }
}
