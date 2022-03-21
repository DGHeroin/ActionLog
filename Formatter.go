package ActionLog

import (
    jsoniter "github.com/json-iterator/go"
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
    data["time"] = entry.Time
    data["msg"] = entry.Message
    var json = jsoniter.ConfigFastest
    bin, err := json.Marshal(&data)
    if err != nil {
       return nil, err
    }
    return append(bin, '\n'), nil
}
