package main

import (
	"github.com/mirogindev/pow-challenge/config"
	"github.com/mirogindev/pow-challenge/internal/db"
	"github.com/mirogindev/pow-challenge/internal/tcpserver"
	"github.com/mirogindev/pow-challenge/internal/timeresolver"
	"github.com/mirogindev/pow-challenge/internal/tools"
	log "github.com/sirupsen/logrus"
	"path"
)

func main() {
	//get current config
	conf := config.GetConfig()
	log.SetLevel(config.GetLogLevelFromString(conf.LogLevel))

	//loading quotes from file
	qu, err := tcpserver.GetQuotesFromFile(path.Join(tools.GetBasePath(), "../../assets", "quotes.txt"))

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("cannot get quotes")
		return
	}

	//init inmemory key/value db
	db, err := db.InitInMemoryDB()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("cannot get quotes")
		return
	}

	ts := tcpserver.TcpServer{
		Host:         conf.Host,
		Port:         conf.Port,
		DB:           db,
		Difficulty:   conf.Difficulty,
		Quotes:       qu,
		TimeResolver: timeresolver.TimeResolverProd{},
	}
	//start server
	err = ts.Start()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Panic("cannot start tcp server")
		return
	}
}
