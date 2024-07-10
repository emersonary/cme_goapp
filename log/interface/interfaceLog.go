package handler

import "github.com/emersonary/go-authentication/log/model"

type LogInterface = interface {
	AddLog(message string, parentLog *model.TLog) *model.TLog
}
