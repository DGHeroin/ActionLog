package ActionLog

import (
    "github.com/sirupsen/logrus"
    "io"
    "io/ioutil"
    "sync"
    "time"
)

type (
    ActionLog struct {
        mutex          sync.Mutex
        log            *logrus.Logger
        pool           sync.Pool
        standardFields logrus.Fields
        sh             *speedyHook
    }
    F map[string]interface{}
)

func New(fields ...F) *ActionLog {
    log := logrus.New()
    log.SetOutput(ioutil.Discard)
    log.SetFormatter(&nullFormatter{})
    sh := &speedyHook{
        w: ioutil.Discard,
    }
    log.AddHook(sh)
    
    L := &ActionLog{
        log: log,
        sh:  sh,
    }
    if len(fields) == 1 {
        var fs map[string]interface{} = fields[0]
        L.standardFields = fs
        sh.standardFields = fs
    }
    L.pool.New = func() interface{} {
        return &logrus.Entry{
            Data: make(logrus.Fields, 6),
        }
    }
    return L
}
func (a *ActionLog) SetWriter(w io.Writer) {
    a.sh.w = w
}
func (a *ActionLog) Info(fields F, args ...interface{}) {
    var ff map[string]interface{} = fields
    
    entry := a.allocEntry()
    
    entry.WithTime(time.Now()).WithFields(ff).WithFields(ff).Info(args...)
    
    a.freeEntry(entry)
}
func (a *ActionLog) allocEntry() *logrus.Entry {
    entry := a.pool.Get().(*logrus.Entry)
    entry.Logger = a.log
    return entry
}

func (a *ActionLog) freeEntry(entry *logrus.Entry) {
    entry.Data = map[string]interface{}{}
    a.pool.Put(entry)
}
