package main

import (
    "fmt"
    "github.com/thedevsaddam/gojsonq/v2"
)

func main()  {
    jq := gojsonq.New().File("./202203-26_12-46-39-+0800-202203-26_12-46-41-+0800_1000000.json")
    var res = jq.Where("code", ">", 999990).
        OrWhere("code", "=", 0).
        Get()

    jq.Reset()
    var arr[]interface{} = res.([]interface{})
    var objs[] map[string]interface{}
    for _, p := range arr {
        objs = append(objs, p.(map[string]interface{}))
    }
    fmt.Println(len(objs))
}
