package tcpserver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mirogindev/pow-challenge/internal/db"
	"github.com/mirogindev/pow-challenge/internal/hashcash"
	"github.com/mirogindev/pow-challenge/internal/protocol"
	"github.com/mirogindev/pow-challenge/internal/timeresolver"
	"github.com/mirogindev/pow-challenge/internal/tools"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTcpServer(t *testing.T) {
	_db, _ := db.InitInMemoryDB()

	difficulty := 3
	remoteClient := "127.0.0.1:58520"
	testTime := time.Now()
	testQuote := "Test quote"
	tr := &timeresolver.TimeResolverMock{
		Time: testTime,
	}

	//Add test resource to db
	col, _ := _db.GetCollection(db.ResourcesCollection)
	col.AddKeyValue(remoteClient, remoteClient)

	ts := TcpServer{
		DB:           _db,
		Difficulty:   difficulty,
		Quotes:       Quotes{testQuote},
		TimeResolver: tr,
	}

	t.Run("Test prepare challenge", func(t *testing.T) {
		ms := &protocol.Message{
			Header: protocol.RequestChallenge,
		}
		res, err := ts.processMessage(ms.Encode(), remoteClient)
		assert.Empty(t, err, "Error while request challenge")
		assert.Equal(t, res.Header, protocol.ResponseChallenge)
		assert.NotEmpty(t, res.Payload)
		hc := &hashcash.HashCashData{}
		json.Unmarshal(res.Payload, hc)

		assert.Equal(t, hc.Date, tools.GetFormattedDateTime(testTime))
		assert.Equal(t, hc.Bits, difficulty)
		assert.Equal(t, hc.Resource, remoteClient)
		assert.Equal(t, hc.Counter, base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", 0))))
	})

	t.Run("Test get quote with valid message", func(t *testing.T) {
		payload := hashcash.HashCashData{
			Ver:      1,
			Bits:     3,
			Date:     "220421192952",
			Resource: "127.0.0.1:58520",
			Rand:     "MTg4ODMy",
			Counter:  "OTkyMzM1",
		}

		j, err := json.Marshal(payload)

		ms := &protocol.Message{
			Header:  protocol.RequestData,
			Payload: j,
		}

		res, err := ts.processMessage(ms.Encode(), remoteClient)
		assert.Empty(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, res.Header, protocol.ResponseData)
		assert.Equal(t, res.Payload, []byte(testQuote))
	})

	t.Run("Test get quote with invalid counter", func(t *testing.T) {
		payload := hashcash.HashCashData{
			Ver:      1,
			Bits:     3,
			Date:     "220421192952",
			Resource: "127.0.0.1:58520",
			Rand:     "MTg4ODMy",
			Counter:  base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", 55))),
		}

		j, err := json.Marshal(payload)

		ms := &protocol.Message{
			Header:  protocol.RequestData,
			Payload: j,
		}

		_, err = ts.processMessage(ms.Encode(), remoteClient)
		assert.NotEmpty(t, err)
		assert.True(t, errors.Is(err, protocol.InvalidChallengeErr))
	})

	t.Run("Test get quote with invalid date", func(t *testing.T) {
		date, _ := tools.ParseDateTime("220421192952")
		date = date.Add(-60 * time.Hour)

		payload := hashcash.HashCashData{
			Ver:      1,
			Bits:     3,
			Date:     tools.GetFormattedDateTime(date),
			Resource: "127.0.0.1:58520",
			Rand:     "MTg4ODMy",
			Counter:  base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v", 55))),
		}

		j, err := json.Marshal(payload)

		ms := &protocol.Message{
			Header:  protocol.RequestData,
			Payload: j,
		}

		_, err = ts.processMessage(ms.Encode(), remoteClient)
		assert.NotEmpty(t, err)
		assert.True(t, errors.Is(err, protocol.DateExprErr))
	})

	t.Run("Test get quote with invalid resource", func(t *testing.T) {

		payload := hashcash.HashCashData{
			Ver:      1,
			Bits:     3,
			Date:     "220421192952",
			Resource: "127.0.0.1:9090",
			Rand:     "MTg4ODMy",
			Counter:  "OTkyMzM1",
		}

		j, err := json.Marshal(payload)

		ms := &protocol.Message{
			Header:  protocol.RequestData,
			Payload: j,
		}

		_, err = ts.processMessage(ms.Encode(), remoteClient)
		assert.NotEmpty(t, err)
		assert.True(t, errors.Is(err, protocol.InvalidResourceErr), err.Error())
	})

}
