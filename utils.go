package ActionLog

import (
    "strconv"
)

var (
    humanSizeSuffixes = [5]string{
        "B", "KB", "MB", "GB", "TB",
    }
)

func HumanFileSize(size float64) string {
    i := 0
    dt := size
    for i = 0; i < 5; i++ {
        if dt <= 1024 {
            break
        }
        dt = dt / 1024
    }
    return strconv.FormatFloat(dt, 'f', 2, 64) +humanSizeSuffixes[i]
}
