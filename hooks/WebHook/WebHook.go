package WebHook

import (
    "bytes"
    "github.com/DGHeroin/ActionLog"
    jsoniter "github.com/json-iterator/go"
    "io/ioutil"
    "log"
    "net/http"
    "sync"
    "time"
)

var (
    RetryPost = 3
)

type (
    webHook struct {
        wg  sync.WaitGroup
        url string
        fn  func(*ActionLog.Entry) bool
        buf *buffer
    }
)

func NewWebHook(interval time.Duration) *webHook {
    h := &webHook{}
    h.buf = NewBuffer(WithFireInterval(interval), WithHandler(func(ts []T) {
        var json = jsoniter.ConfigCompatibleWithStandardLibrary
        bin, err := json.Marshal(ts)
        if err != nil {
            return
        }
        h.doPost(string(bin))
    }))
    return h
}
func (h *webHook) Fire(entry *ActionLog.Entry) error {
    // 过滤失败
    if !h.fn(entry) {
        return nil
    }
    data := make(ActionLog.F, len(entry.Data)+2)
    for k, v := range entry.Data {
        data[k] = v
    }
    data["time"] = entry.Time.String()
    if entry.Message != "" {
        data["msg"] = entry.Message
    }

    h.buf.Add(T(data))
    return nil
}
func (h *webHook) AddHook(url string, fn func(*ActionLog.Entry) bool) {
    h.url = url
    h.fn = fn
}
func (h *webHook) Drain() {
    h.buf.Drain()
    h.wg.Wait()
}
func (h *webHook) doPost(body string) {
    h.wg.Add(1)
    defer func() {
        recover()
        h.wg.Done()
    }()
    // 发送
    retry := 0
    for retry < RetryPost {
        if resp, err := http.Post(h.url, "application/json", bytes.NewBufferString(body)); err == nil {
            data, _ := ioutil.ReadAll(resp.Body)
            log.Println(len(body), string(data))
            break
        }
        time.Sleep(time.Second)
        retry++
    }
}
