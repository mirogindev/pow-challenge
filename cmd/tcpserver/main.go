package main

import (
	"github.com/mirogindev/pow-challenge/internal/db"
	"github.com/mirogindev/pow-challenge/internal/tcpserver"
	log "github.com/sirupsen/logrus"
	"path"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	bp         = filepath.Dir(b)
)

func main() {
	log.SetLevel(log.TraceLevel)
	qu, err := tcpserver.GetQuotesFromFile(path.Join(b, "../../../assets", "quotes.txt"))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("cannot get quotes")
		return
	}

	db, err := db.InitInMemoryDB()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("cannot get quotes")
		return
	}

	ts := tcpserver.TcpServer{
		Port:       8093,
		DB:         db,
		Difficulty: 4,
		Quotes:     qu,
	}
	err = ts.Start()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("cannot start tcp server")
		return
	}
}
