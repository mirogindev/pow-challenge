package tcpclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mirogindev/pow-challenge/internal/hashcash"
	"github.com/mirogindev/pow-challenge/internal/protocol"
	log "github.com/sirupsen/logrus"
	"net"
)

type TcpClient struct {
	MaxIters int64
	Addr     string
	conn     net.Conn
}

func (c *TcpClient) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"addr": c.Addr,
	}).Println("successfully connected to tcp server")
	return nil
}

func (c *TcpClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		log.WithFields(log.Fields{
			"addr": err,
		}).Errorf("cannot close connection")
		return err
	}
	return nil
}

func (c *TcpClient) CloseSession() error {
	_, err := c.conn.Write([]byte(fmt.Sprintf("%s\n", "STOP")))
	if err != nil {
		log.WithFields(log.Fields{
			"addr": err,
		}).Errorf("cannot close session")
		return err
	}
	return nil
}

func (c *TcpClient) SendMessageWithReply(ms *protocol.Message) (*protocol.Message, error) {
	en := ms.Encode()

	err := c.sendMessage(en)

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"msg": en,
	}).Debug("sent message")

	data, _ := bufio.NewReader(c.conn).ReadString('\n')
	p, err := protocol.DecodeMessage(data)

	if err != nil {
		return nil, err
	}

	var hashData hashcash.HashCashData
	err = json.Unmarshal(p.Payload, &hashData)

	log.WithFields(log.Fields{
		"msg": hashData.ToString(),
	}).Debug("challenge received")

	if err != nil {
		return nil, protocol.UnmarshallErr
	}
	resolvedChallenge, err := hashcash.ResolveChallenge(hashData, c.MaxIters)

	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(resolvedChallenge)

	solvedMessage := &protocol.Message{Header: protocol.RequestData, Payload: b}

	err = c.sendMessage(solvedMessage.Encode())

	data, _ = bufio.NewReader(c.conn).ReadString('\n')
	p, err = protocol.DecodeMessage(data)

	log.WithFields(log.Fields{
		"quote": string(p.Payload),
	}).Debug("quote received")

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (c *TcpClient) sendMessage(m string) error {
	_, err := c.conn.Write([]byte(fmt.Sprintf("%s\n", m)))
	return err

}
