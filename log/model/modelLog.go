package model

import (
	"time"

	"github.com/google/uuid"
)

type LogType int

const (
	LT_HTTP LogType = iota
	LT_DB
	LT_CONFIG
	LT_ERROR
)

type TLog struct {
	Id        uuid.UUID
	CreatedAt time.Time
	LogText   string
	LogType   LogType
	ParentLog *TLog
}

func NewLog(logText string, parentLog *TLog) *TLog {
	return &TLog{Id: uuid.Nil, LogText: logText, ParentLog: parentLog}
}

func NewLogEmpty() *TLog {
	return &TLog{Id: uuid.Nil}
}
