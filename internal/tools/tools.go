package tools

import (
	"github.com/mirogindev/pow-challenge/internal/timeresolver"
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

func CheckDateExpired(d time.Time, tr timeresolver.TimeResolver) bool {
	expiresAt := d.Add(48 * time.Hour)
	return tr.Now().Unix() > expiresAt.Unix()
}

func GetBasePath() string {
	return basePath
}
