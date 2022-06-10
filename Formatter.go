package ActionLog

import (
    "bytes"
    jsoniter "github.com/json-iterator/go"
)

type (
    Formatter interface {
        SetPrefix(string)
        Format(entry *Entry) ([]byte, error)
    }
    defaultFormatter struct {
        prefix []byte
    }
)

func (d *defaultFormatter) SetPrefix(s string) {
    d.prefix = []byte(s)
}

func (d defaultFormatter) Format(entry *Entry) ([]byte, error) {
    data := make(F, len(entry.Data)+2)
    for k, v := range entry.Data {
        data[k] = v
    }
    data["time"] = entry.Time.String()
    if entry.Message != "" {
        data["msg"] = entry.Message
    }
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    bin, err := json.Marshal(&data)
    if err != nil {
        return nil, err
    }
    if d.prefix == nil {
        return append(bin, '\n'), nil
    } else {
        buf := bytes.NewBuffer(d.prefix)
        buf.Write(bin)
        buf.WriteString("\n")
        return buf.Bytes(), nil
    }
}
