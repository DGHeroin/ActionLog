package main

import (
    "github.com/DGHeroin/ActionLog"
    "os"
)

func main() {
    L := ActionLog.New(ActionLog.F{
        "hostname": "staging-1",
        "appname":  "ice-sword",
        "session":  "1c3b3r9",
    })
    L.SetWriter(os.Stdout)
    L.Info(ActionLog.F{"code": 10, "playerId": "aaa"}, "hello \n", "world")
    L.Info(ActionLog.F{"code": 11, "playerId": "bbb"})
    
}
