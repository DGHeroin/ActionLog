package ActionLog

import (
    jsoniter "github.com/json-iterator/go"
    "github.com/sirupsen/logrus"
    "io"
    "time"
)

type (
    speedyHook struct {
        w              io.Writer
        standardFields logrus.Fields
    }
    nullFormatter struct {
    }
)

func (s *speedyHook) Levels() []logrus.Level {
    return []logrus.Level{logrus.InfoLevel}
}

func (s *speedyHook) Fire(entry *logrus.Entry) error {
    data := make(F, len(entry.Data)+2+len(s.standardFields))
    for k, v := range entry.Data {
        data[k] = v
        //switch v := v.(type) {
        //case error:
        //    // Otherwise errors are ignored by `encoding/json`
        //    // https://github.com/sirupsen/logrus/issues/137
        //    data[k] = v.Error()
        //default:
        //    data[k] = v
        //}
    }
    for k, v := range s.standardFields {
        data[k] = v
    }
    data["msg"] = entry.Message
    data["time"] = time.Now()
    var json = jsoniter.ConfigCompatibleWithStandardLibrary
    jsd, _ := json.Marshal(&data)
    _, err := s.w.Write(append(jsd, []byte("\n")...))
    return err
}

func (f *nullFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    return []byte(""), nil
}
