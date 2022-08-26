package main

import (
    "github.com/DGHeroin/ActionLog"
    "github.com/DGHeroin/ActionLog/hooks/WebHook"
    "os"
    "time"
)

func main() {
    L := ActionLog.New(ActionLog.F{
        "hostname": "staging-1",
        "appname":  "ice-sword",
        "session":  "1c3b3r9",
    })
    L.SetWriter(os.Stdout)
    hook := WebHook.NewWebHook(time.Second)
    hook.AddHook("https://httpbin.org/post", func(entry *ActionLog.Entry) bool {
        if entry.Data["type"] == "login" {
            return true
        }
        return false
    })
    L.AddHook(hook)
    L.Info(ActionLog.F{"code": 10, "playerId": "aaa", "type": "login"}, "hello \n", "world")
    L.Info(ActionLog.F{"code": 11, "playerId": "bbb", "type": "logout"})
    L.Info(ActionLog.F{"code": 12, "playerId": "ccc", "type": "login"}, "hello \n", "world")
    time.Sleep(time.Second)
    L.Info(ActionLog.F{"code": 66, "playerId": "ddd", "type": "login"}, "hello \n", "world")
    hook.Drain()
}
