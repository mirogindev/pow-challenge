package hashcash

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/mirogindev/pow-challenge/internal/timeresolver"
	"github.com/mirogindev/pow-challenge/internal/tools"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

const (
	Zero rune = 48
)

var (
	LimitExceedErr = errors.New("hash not resolved, iterations limit exceed")
)

type HashCashData struct {
	Ver      int
	Bits     int
	Date     string
	Resource string
	Rand     string
	Counter  string
}

func (d *HashCashData) ToString() string {
	return fmt.Sprintf("%v:%d:%s:%s::%s:%s", d.Ver, d.Bits, d.Date, d.Resource, d.Rand, d.Counter)
}

// GenerateChallenge - generate challenge for client
func GenerateChallenge(remoteClient string, difficulty int, ts timeresolver.TimeResolver) *HashCashData {

	//generate a random number for challenge
	rand := rand.Intn(1000000)
	ms := &HashCashData{
		Ver:      1,
		Bits:     difficulty,
		Rand:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", rand))),
		Counter:  base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", 0))),
		Date:     tools.GetFormattedDateTime(ts.Now()),
		Resource: remoteClient,
	}

	return ms
}

// ResolveChallenge - method which trying to resolve challenge
func ResolveChallenge(m HashCashData, maxIters int64) (*HashCashData, error) {
	var counter int64

	var initRanVal = rand.Intn(1000000)
	for counter <= maxIters {
		m.Counter = base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", initRanVal)))
		if h, ok := IsChallengeValid(&m); ok {
			log.WithFields(log.Fields{
				"msg":        m.ToString(),
				"hash":       h,
				"iterations": counter,
			}).Println("challenge resolved")
			return &m, nil
		} else {
			counter++
			initRanVal++
		}
	}
	return nil, LimitExceedErr
}

func IsChallengeValid(m *HashCashData) (string, bool) {
	ms := m.ToString()
	sha := getSha1(ms)
	valid := checkZeros(sha, Zero, m.Bits)
	return sha, valid
}

func getSha1(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func checkZeros(hash string, char rune, n int) bool {
	for _, val := range hash[:n] {
		if val != char {
			return false
		}
	}
	return true
}
