package protocol

import "errors"

const (
	End               = "End"
	RequestChallenge  = "RequestChallenge"
	ResponseChallenge = "ResponseChallenge"
	RequestData       = "RequestData"
	ResponseData      = "ResponseData"
)

var (
	QuitErr             = errors.New("session closed")
	SerErr              = errors.New("serialization error")
	InvalidComErr       = errors.New("invalid protocol command")
	InvalidChallengeErr = errors.New("invalid challenge command")
	UnmarshallErr       = errors.New("cannot unmarshall hashdata")
	ParseDateErr        = errors.New("cannot parse datetime")
	DateExprErr         = errors.New("message header older than two days")
	InvalidResourceErr  = errors.New("resource not found")
	HashExitsErr        = errors.New("hash not found")
)
