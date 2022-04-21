package tools

import "time"

const layout = "060102150405"

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
