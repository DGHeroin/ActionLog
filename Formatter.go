package ActionLog

import (
    jsoniter "github.com/json-iterator/go"
    "time"
)

type (
    Formatter interface {
        Format(entry *Entry) ([]byte, error)
    }
    defaultFormatter struct {
    }
)

func (d defaultFormatter) Format(entry *Entry) ([]byte, error) {
    data := make(F, len(entry.Data)+2)
    for k, v := range entry.Data {
        data[k] = v
    }
    data["msg"] = entry.Message
    data["time"] = time.Now()
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    bin, err := json.Marshal(&data)
    if err != nil {
        return nil, err
    }
    return bin, nil
}
