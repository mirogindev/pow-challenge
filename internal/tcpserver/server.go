package tcpserver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mirogindev/pow-challenge/internal/db"
	"github.com/mirogindev/pow-challenge/internal/hashcash"
	"github.com/mirogindev/pow-challenge/internal/protocol"
	"github.com/mirogindev/pow-challenge/internal/timeresolver"
	"github.com/mirogindev/pow-challenge/internal/tools"
	log "github.com/sirupsen/logrus"
	"net"
)

type TcpServer struct {
	Host         string
	Port         int
	Difficulty   int
	Quotes       Quotes
	DB           db.DB
	TimeResolver timeresolver.TimeResolver
}

func (s *TcpServer) Start() error {
	l, err := net.Listen("tcp4", fmt.Sprintf(":%v", s.Port))

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"host":       s.Host,
		"port":       s.Port,
		"difficulty": s.Difficulty,
		"quotes":     len(s.Quotes),
	}).Println("server started successfully")

	defer l.Close()

	for {

		c, err := l.Accept()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Errorf("cannot process connection")
			c.Close()
			continue
		}
		go s.handleConnection(c)
	}
}

//Handle new tcp connection
func (s *TcpServer) handleConnection(c net.Conn) {
	log.WithFields(log.Fields{
		"remoteAddr": c.RemoteAddr().String(),
	}).Println("process connection")

	defer c.Close()

	for {
		remoteClient := c.RemoteAddr().String()
		packet, err := bufio.NewReader(c).ReadString('\n')

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("cannot parse packet")
			break
		}

		resp, err := s.processMessage(packet, remoteClient)

		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("session stopped by message processor")
			break
		}

		_, err = c.Write([]byte(fmt.Sprintf("%s\n", resp.Encode())))
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("cannot respond to a client")
		}
	}
}

//This method processed all clients requests
func (s *TcpServer) processMessage(data, remoteClient string) (*protocol.Message, error) {
	hashCol, _ := s.DB.GetCollection(db.HashCollection)
	resCol, _ := s.DB.GetCollection(db.ResourcesCollection)
	p, err := protocol.DecodeMessage(data)
	if err != nil {
		return nil, err
	}

	switch p.Header {
	case protocol.End:

		return nil, protocol.QuitErr
	case protocol.RequestData:
		var hashData hashcash.HashCashData

		err := json.Unmarshal(p.Payload, &hashData)
		requestLogger := log.WithFields(log.Fields{"msg": hashData, "resource": remoteClient})
		if err != nil {
			requestLogger.Traceln(protocol.UnmarshallErr)
			return nil, protocol.UnmarshallErr
		}
		//Check if remote client exits in the database
		if !resCol.KeyExist(hashData.Resource) {
			requestLogger.Traceln(protocol.InvalidResourceErr)
			return nil, protocol.InvalidResourceErr
		}

		hd, err := tools.ParseDateTime(hashData.Date)
		if err != nil {
			requestLogger.Traceln(protocol.ParseDateErr)
			return nil, protocol.ParseDateErr
		}
		//Check if date generate no longer than 2 days ago
		//as described here https://ru.wikipedia.org/wiki/Hashcash
		exp := tools.CheckDateExpired(hd, s.TimeResolver)
		if exp {
			requestLogger.Traceln(protocol.DateExprErr)
			return nil, protocol.DateExprErr
		} else {
			requestLogger.Traceln("received a valid date")
		}
		//Check if hash has required number of zeros
		if hash, ok := hashcash.IsChallengeValid(&hashData); ok {

			//Checking that the hash has not been used
			// if not  add it to the db
			if hashCol.KeyExist(hash) {
				log.WithFields(log.Fields{
					"msg": hashData.ToString(),
				}).Error(protocol.HashExitsErr)
				return nil, protocol.HashExitsErr
			} else {
				hashCol.AddKeyValue(hash, hash)
			}

			log.WithFields(log.Fields{
				"msg": hashData.ToString(),
			}).Debug("received a valid challenge")

			//If all test passed  then send to a client a random quote
			ms := &protocol.Message{Header: protocol.ResponseData, Payload: []byte(s.Quotes.GetRandomQuote())}
			return ms, nil
		} else {
			return nil, protocol.InvalidChallengeErr
		}

	case protocol.RequestChallenge:
		//Generate challenge for the client
		hd := hashcash.GenerateChallenge(remoteClient, s.Difficulty, s.TimeResolver)

		//Add remote client to the db
		if !resCol.KeyExist(remoteClient) {
			err := resCol.AddKeyValue(remoteClient, remoteClient)
			if err != nil {
				return nil, err
			}
		}

		payload, err := json.Marshal(hd)

		if err != nil {
			log.Println(err)
			return nil, protocol.SerErr
		}

		//Send this challenge back to the client
		ms := &protocol.Message{Header: protocol.ResponseChallenge, Payload: payload}
		return ms, nil
	}

	return nil, protocol.InvalidComErr
}
