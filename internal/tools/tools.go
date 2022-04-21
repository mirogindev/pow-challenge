package tools

import (
	"path/filepath"
	"runtime"
	"time"
)

const layout = "060102150405"

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

func GetFormattedDateTime(t time.Time) string {
	return t.Format(layout)
}

func ParseDateTime(s string) (time.Time, error) {
	return time.Parse(layout, s)
}

func CheckDateExpired(d time.Time) bool {
	expiresAt := d.Add(48 * time.Hour)
	return time.Now().Unix() > expiresAt.Unix()
}

func GetBasePath() string {
	return basePath
}
