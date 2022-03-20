package ActionLog

import (
    "bytes"
    "fmt"
    "time"
)

type (
    Entry struct {
        Data    F
        Time    time.Time
        Message string
        Buffer  *bytes.Buffer
    }
)

func (e *Entry) WithTime(t time.Time) *Entry {
    e.Time = t
    return e
}
func (e *Entry) WithField(key string, value interface{}) *Entry {
    e.Data[key] = value
    return e
}
func (e *Entry) WithFields(fields F) *Entry {
    for k, v := range fields {
        e.Data[k] = v
    }
    return e
}
func (e *Entry) Info(args ...interface{}) {
    if len(args) == 0 {
        return
    }
    e.Message = fmt.Sprint(args)
}
