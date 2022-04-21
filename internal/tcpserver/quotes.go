package tcpserver

import (
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

type Quotes []string

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (q Quotes) GetRandomQuote() string {
	return q[rand.Intn(len(q))]
}

func GetQuotesFromFile(p string) (Quotes, error) {
	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(content), "\n"), nil
}
