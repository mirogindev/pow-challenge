package main

import (
	"errors"
	"fmt"
	"github.com/mirogindev/pow-challenge/config"
	"github.com/mirogindev/pow-challenge/internal/hashcash"
	"github.com/mirogindev/pow-challenge/internal/protocol"
	"github.com/mirogindev/pow-challenge/internal/tcpclient"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	log.SetLevel(log.TraceLevel)
	conf := config.GetConfig()
	tc := &tcpclient.TcpClient{Addr: fmt.Sprintf("%s:%v", conf.Host, conf.Port), MaxIters: conf.MaxIterations}
	err := tc.Connect()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatalf("server unavailable")
	}

	defer tc.CloseSession()

	for {
		time.Sleep(1 * time.Second)
		ms := &protocol.Message{
			Header: protocol.RequestChallenge,
		}

		_, err = tc.SendMessageWithReply(ms)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Errorln("cannot process message")
			if errors.Is(err, hashcash.LimitExceedErr) {
				continue
			} else {
				log.WithFields(log.Fields{
					"error": err,
				}).Fatalf("something went wrong, stoping client")
			}
		}
	}
}
