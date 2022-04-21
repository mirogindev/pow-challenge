package main

import (
	"github.com/mirogindev/pow-challenge/config"
	"github.com/mirogindev/pow-challenge/internal/db"
	"github.com/mirogindev/pow-challenge/internal/tcpserver"
	"github.com/mirogindev/pow-challenge/internal/tools"
	log "github.com/sirupsen/logrus"
	"path"
)

func main() {
	log.SetLevel(log.TraceLevel)
	qu, err := tcpserver.GetQuotesFromFile(path.Join(tools.GetBasePath(), "../../assets", "quotes.txt"))

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

	conf := config.GetConfig()

	ts := tcpserver.TcpServer{
		Host:       conf.Host,
		Port:       conf.Port,
		DB:         db,
		Difficulty: conf.Difficulty,
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
